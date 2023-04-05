package registry

import (
	"go/types"
	"strconv"
)

// MethodScope is the sub-registry for allocating variables present in
// the method scope.
//
// It should be created using a registry instance.
type MethodScope struct {
	registry *Registry

	vars       []*Var
	conflicted map[string]bool
}

// populateImports extracts all the package imports for a given type
// recursively. The imported packages by a single type can be more than
// one (ex: map[a.Type]b.Type).
func (m MethodScope) populateImports(t types.Type, imports map[string]*Package) {
	switch t := t.(type) {
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			imports[StripVendorPath(pkg.Path())] = m.registry.AddImport(pkg)
		}
	}
}

func (m *MethodScope) AddVar(vr *types.Var, suffix string) *Var {
	imports := make(map[string]*Package)
	m.populateImports(vr.Type(), imports)
	m.resolveImportVarConflicts(imports)

	name := varName(vr, suffix)
	// Ensure that the var name does not conflict with a package import.
	if _, ok := m.registry.searchImport(name); ok {
		name += "MoqParam"
	}
	if _, ok := m.searchVar(name); ok || m.conflicted[name] {
		name = m.resolveVarNameConflict(name)
	}

	v := Var{
		vr:      vr,
		imports: imports,
		Name:    name,
	}
	m.vars = append(m.vars, &v)
	return &v
}

func varName(vr *types.Var, suffix string) string {
	name := vr.Name()
	if name != "" && name != "_" {
		return name + suffix
	}

	name = varNameForType(vr.Type()) + suffix

	switch name {
	case "mock", "callInfo", "break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for",
		"import", "return", "var",
		// avoid shadowing basic types
		"string", "bool", "byte", "rune", "uintptr",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "complex64", "complex128":
		name += "MoqParam"
	}

	return name
}

// resolveImportVarConflicts ensures that all the newly added imports do not
// conflict with any of the existing vars.
func (m MethodScope) resolveImportVarConflicts(imports map[string]*Package) {
	// Ensure that all the newly added imports do not conflict with any of the
	// existing vars.
	for _, imprt := range imports {
		if v, ok := m.searchVar(imprt.Qualifier()); ok {
			v.Name += "MoqParam"
		}
	}
}

func (m MethodScope) searchVar(name string) (*Var, bool) {
	for _, v := range m.vars {
		if v.Name == name {
			return v, true
		}
	}

	return nil, false
}
func (m *MethodScope) resolveVarNameConflict(suggested string) string {
	for n := 1; ; n++ {
		_, ok := m.searchVar(suggested + strconv.Itoa(n))
		if ok {
			continue
		}

		if n == 1 {
			conflict, _ := m.searchVar(suggested)
			conflict.Name += "1"
			m.conflicted[suggested] = true
			n++
		}
		return suggested + strconv.Itoa(n)
	}
}
