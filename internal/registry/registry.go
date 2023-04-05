package registry

import (
	"errors"
	"fmt"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type Registry struct {
	srcPkg  *packages.Package
	aliases map[string]string
	imports map[string]*Package
}

// SrcPkg returns the types info for the source package.
func (r Registry) SrcPkg() *types.Package {
	return r.srcPkg.Types
}

// New loads the source package info and returns a new instance of
// Registry.
func New(srcDir string) (*Registry, error) {
	srcPkg, err := pkgInfoFromPath(
		srcDir, packages.NeedName|packages.NeedSyntax|packages.NeedTypes|packages.NeedTypesInfo|packages.NeedDeps,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't load source package: %s", err)
	}

	return &Registry{
		srcPkg:  srcPkg,
		imports: make(map[string]*Package),
	}, nil
}

type Var struct {
	vr      *types.Var
	imports map[string]*Package

	Name string
}

func findPkgPath(pkgInputVal string, srcPkg *packages.Package) string {
	if pkgInputVal == "" {
		return srcPkg.PkgPath
	}
	if pkgInDir(srcPkg.PkgPath, pkgInputVal) {
		return srcPkg.PkgPath
	}
	subdirectoryPath := filepath.Join(srcPkg.PkgPath, pkgInputVal)
	if pkgInDir(subdirectoryPath, pkgInputVal) {
		return subdirectoryPath
	}
	return ""
}

func pkgInDir(pkgName, dir string) bool {
	currentPkg, err := pkgInfoFromPath(dir, packages.NeedName)
	if err != nil {
		return false
	}
	return currentPkg.Name == pkgName || currentPkg.Name+"_test" == pkgName
}

func (r Registry) LookupStruct(name string) (*types.Struct, *types.TypeParamList, error) {
	srcPkg := r.SrcPkg()
	scope := srcPkg.Scope()
	obj := scope.Lookup(name)
	if obj == nil {
		return nil, nil, fmt.Errorf("struct not found: %s", name)
	}

	var tparams *types.TypeParamList
	named, ok := obj.Type().(*types.Named)
	if ok {
		tparams = named.TypeParams()
	}

	return obj.Type().Underlying().(*types.Struct), tparams, nil
}

func pkgInfoFromPath(srcDir string, mode packages.LoadMode) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: mode,
		Dir:  srcDir,
	})
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("package not found")
	}
	if len(pkgs) > 1 {
		return nil, errors.New("found more than one package")
	}
	if errs := pkgs[0].Errors; len(errs) != 0 {
		if len(errs) == 1 {
			return nil, errs[0]
		}
		return nil, fmt.Errorf("%s (and %d more errors)", errs[0], len(errs)-1)
	}
	return pkgs[0], nil
}

func (r Registry) SrcPkgName() string {
	return r.srcPkg.Name
}

// AddImport adds the given package to the set of imports. It generates a
// suitable alias if there are any conflicts with previously imported
// packages.
func (r *Registry) AddImport(pkg *types.Package) *Package {
	path := StripVendorPath(pkg.Path())

	if imprt, ok := r.imports[path]; ok {
		return imprt
	}

	imprt := Package{pkg: pkg, Alias: r.aliases[path]}

	if conflict, ok := r.searchImport(imprt.Qualifier()); ok {
		r.resolveImportConflict(&imprt, conflict, 0)
	}

	r.imports[path] = &imprt
	return &imprt
}

func (r Registry) searchImport(name string) (*Package, bool) {
	for _, imprt := range r.imports {
		if imprt.Qualifier() == name {
			return imprt, true
		}
	}

	return nil, false
}

// MethodScope returns a new MethodScope.
func (r *Registry) MethodScope() *MethodScope {
	return &MethodScope{
		registry:   r,
		conflicted: map[string]bool{},
	}
}

// resolveImportConflict generates and assigns a unique alias for
// packages with conflicting qualifiers.
func (r Registry) resolveImportConflict(a, b *Package, lvl int) {
	if a.uniqueName(lvl) == b.uniqueName(lvl) {
		r.resolveImportConflict(a, b, lvl+1)
		return
	}

	for _, p := range []*Package{a, b} {
		name := p.uniqueName(lvl)
		// Even though the name is not conflicting with the other package we
		// got, the new name we want to pick might already be taken. So check
		// again for conflicts and resolve them as well. Since the name for
		// this package would also get set in the recursive function call, skip
		// setting the alias after it.
		if conflict, ok := r.searchImport(name); ok && conflict != p {
			r.resolveImportConflict(p, conflict, lvl+1)
			continue
		}

		p.Alias = name
	}
}
