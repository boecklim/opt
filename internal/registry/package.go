package registry

import (
	"go/types"
	"path"
	"sort"
	"strings"
)

type Package struct {
	pkg *types.Package

	Alias string
}

// StripVendorPath strips the vendor dir prefix from a package path.
// For example we might encounter an absolute path like
// github.com/foo/bar/vendor/github.com/pkg/errors which is resolved
// to github.com/pkg/errors.
func StripVendorPath(p string) string {
	parts := strings.Split(p, "/vendor/")
	if len(parts) == 1 {
		return p
	}
	return strings.TrimLeft(path.Join(parts[1:]...), "/")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func reverse(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// Qualifier returns the qualifier which must be used to refer to types
// declared in the package.
func (p *Package) Qualifier() string {
	if p == nil {
		return ""
	}

	if p.Alias != "" {
		return p.Alias
	}

	return p.pkg.Name()
}

var replacer = strings.NewReplacer(
	"go-", "",
	"-go", "",
	"-", "",
	"_", "",
	".", "",
	"@", "",
	"+", "",
	"~", "",
)

// Imports returns the list of imported packages. The list is sorted by
// path.
func (r Registry) Imports() []*Package {
	imports := make([]*Package, 0, len(r.imports))
	for _, imprt := range r.imports {
		imports = append(imports, imprt)
	}
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path() < imports[j].Path()
	})
	return imports
}

// uniqueName generates a unique name for a package by concatenating
// path components. The generated name is guaranteed to unique with an
// appropriate level because the full package import paths themselves
// are unique.
func (p Package) uniqueName(lvl int) string {
	pp := strings.Split(p.Path(), "/")
	reverse(pp)

	var name string
	for i := 0; i < min(len(pp), lvl+1); i++ {
		name = strings.ToLower(replacer.Replace(pp[i])) + name
	}

	return name
}

// Path is the full package import path (without vendor).
func (p *Package) Path() string {
	if p == nil {
		return ""
	}

	return StripVendorPath(p.pkg.Path())
}
