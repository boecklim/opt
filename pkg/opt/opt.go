package opt

import (
	"bytes"
	"errors"
	"go/format"
	"go/types"
	"io"
	"opt/internal/registry"

	"opt/internal/template"
)

type Generator struct {
	cfg Config

	registry *registry.Registry
	tmpl     template.Template
}

// Config specifies details about how interfaces should be mocked.
// SrcDir is the only field which needs be specified.
type Config struct {
	SrcDir string
	// Formatter  string
	// StubImpl   bool
	// SkipEnsure bool
	// WithResets bool
}

func New(cfg Config) (*Generator, error) {
	reg, err := registry.New(cfg.SrcDir)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	return &Generator{
		cfg:      cfg,
		registry: reg,
		tmpl:     tmpl,
	}, nil
}

func (m *Generator) Generate(w io.Writer, structName string) error {
	if len(structName) == 0 {
		return errors.New("must specify one struct")
	}

	structType, _, err := m.registry.LookupStruct(structName)
	if err != nil {
		return err
	}

	members := make([]template.Member, structType.NumFields())
	imports := map[string]*registry.Package{}

	for i := 0; i < structType.NumFields(); i++ {

		m.populateImports(structType.Field(i).Type(), imports)
		// switch t := structType.Field(i).Type().(type) {
		// case *types.Named:
		// 	if pkg := t.Obj().Pkg(); pkg != nil {
		// 		m.registry.AddImport(pkg)
		// 	}
		// }

		members[i] = template.Member{
			Name:       structType.Field(i).Name(),
			Type:       structType.Field(i).Type().String(),
			StructName: structName,
			// TypeParams:    m.typeParams(tparams),
		}
	}

	data := template.Data{
		PkgName:    m.registry.SrcPkgName(),
		StructName: structName,
		Members:    members,
	}

	// imprt := m.registry.AddImport(m.registry.SrcPkg())
	// data.SrcPkgQualifier = imprt.Qualifier() + "."

	data.Imports = m.registry.Imports()
	var buf bytes.Buffer
	err = m.tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(formatted)
	if err != nil {
		return err
	}
	return nil
}
func (m *Generator) populateImports(t types.Type, imports map[string]*registry.Package) {
	switch t := t.(type) {
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			imports[registry.StripVendorPath(pkg.Path())] = m.registry.AddImport(pkg)
		}

	case *types.Array:
		m.populateImports(t.Elem(), imports)

	case *types.Slice:
		m.populateImports(t.Elem(), imports)

	case *types.Signature:
		for i := 0; i < t.Params().Len(); i++ {
			m.populateImports(t.Params().At(i).Type(), imports)
		}
		for i := 0; i < t.Results().Len(); i++ {
			m.populateImports(t.Results().At(i).Type(), imports)
		}

	case *types.Map:
		m.populateImports(t.Key(), imports)
		m.populateImports(t.Elem(), imports)

	case *types.Chan:
		m.populateImports(t.Elem(), imports)

	case *types.Pointer:
		m.populateImports(t.Elem(), imports)

	case *types.Struct: // anonymous struct
		for i := 0; i < t.NumFields(); i++ {
			m.populateImports(t.Field(i).Type(), imports)
		}

	case *types.Interface: // anonymous interface
		for i := 0; i < t.NumExplicitMethods(); i++ {
			m.populateImports(t.ExplicitMethod(i).Type(), imports)
		}
		for i := 0; i < t.NumEmbeddeds(); i++ {
			m.populateImports(t.EmbeddedType(i), imports)
		}
	}
}
