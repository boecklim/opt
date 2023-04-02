package registry

import (
	"errors"
	"fmt"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type Registry struct {
	srcPkg     *packages.Package
	moqPkgPath string
}

// SrcPkg returns the types info for the source package.
func (r Registry) SrcPkg() *types.Package {
	return r.srcPkg.Types
}

// New loads the source package info and returns a new instance of
// Registry.
func New(srcDir, moqPkg string) (*Registry, error) {
	srcPkg, err := pkgInfoFromPath(
		srcDir, packages.NeedName|packages.NeedSyntax|packages.NeedTypes|packages.NeedTypesInfo|packages.NeedDeps,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't load source package: %s", err)
	}

	return &Registry{
		srcPkg:     srcPkg,
		moqPkgPath: findPkgPath(moqPkg, srcPkg),
	}, nil
}

type Var struct {
	vr         *types.Var
	imports    map[string]*Package
	moqPkgPath string

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

func (r Registry) LookupInterface(name string) (*types.Interface, *types.TypeParamList, error) {
	obj := r.SrcPkg().Scope().Lookup(name)
	if obj == nil {
		return nil, nil, fmt.Errorf("interface not found: %s", name)
	}

	if !types.IsInterface(obj.Type()) {
		return nil, nil, fmt.Errorf("%s (%s) is not an interface", name, obj.Type())
	}

	var tparams *types.TypeParamList
	named, ok := obj.Type().(*types.Named)
	if ok {
		tparams = named.TypeParams()
	}

	return obj.Type().Underlying().(*types.Interface).Complete(), tparams, nil
}

func (r Registry) LookupStruct(name string) (*types.Struct, *types.TypeParamList, error) {
	srcPkg := r.SrcPkg()
	scope := srcPkg.Scope()
	obj := scope.Lookup(name)
	if obj == nil {
		return nil, nil, fmt.Errorf("interface not found: %s", name)
	}

	// if !types.IsInterface(obj.Type()) {
	// 	return nil, nil, fmt.Errorf("%s (%s) is not an interface", name, obj.Type())
	// }

	var tparams *types.TypeParamList
	named, ok := obj.Type().(*types.Named)
	if ok {
		tparams = named.TypeParams()
	}

	return obj.Type().Underlying().(*types.Struct), tparams, nil

	//obj.Type().Underlying().(*types.Struct).Field(), tparams, nil
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
