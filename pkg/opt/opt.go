package opt

import (
	"bytes"
	"errors"
	"go/format"
	"go/types"
	"io"
	"opt/internal/registry"
	"strings"

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

	optStart := false

	members := make([]template.Member, 0, structType.NumFields())
	parameterMembers := make([]template.Member, 0, structType.NumFields())

	for i := 0; i < structType.NumFields(); i++ {

		m.populateImports(structType.Field(i).Type())

		newMember := template.Member{
			Name:            structType.Field(i).Name(),
			CapitalizedName: capitalise(structType.Field(i).Name()),
			Type:            structType.Field(i).Type().String(),
			StructName:      structName,
		}

		tag := structType.Tag(i)
		if strings.Contains(tag, "opt:\"true\"") {
			optStart = true
		}

		if optStart {
			members = append(members, newMember)
		} else {
			parameterMembers = append(parameterMembers, newMember)
		}

	}

	data := template.Data{
		PkgName:          m.registry.SrcPkgName(),
		StructName:       structName,
		Members:          members,
		ParameterMembers: parameterMembers,
	}

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
func (m *Generator) populateImports(t types.Type) {
	switch t := t.(type) {
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			m.registry.AddImport(pkg)
		}

	case *types.Array:
		m.populateImports(t.Elem())

	case *types.Slice:
		m.populateImports(t.Elem())

	case *types.Signature:
		for i := 0; i < t.Params().Len(); i++ {
			m.populateImports(t.Params().At(i).Type())
		}
		for i := 0; i < t.Results().Len(); i++ {
			m.populateImports(t.Results().At(i).Type())
		}

	case *types.Map:
		m.populateImports(t.Key())
		m.populateImports(t.Elem())

	case *types.Chan:
		m.populateImports(t.Elem())

	case *types.Pointer:
		m.populateImports(t.Elem())

	case *types.Struct: // anonymous struct
		for i := 0; i < t.NumFields(); i++ {
			m.populateImports(t.Field(i).Type())
		}

	case *types.Interface: // anonymous interface
		for i := 0; i < t.NumExplicitMethods(); i++ {
			m.populateImports(t.ExplicitMethod(i).Type())
		}
		for i := 0; i < t.NumEmbeddeds(); i++ {
			m.populateImports(t.EmbeddedType(i))
		}
	}
}

func capitalise(s string) string { return strings.ToUpper(s[:1]) + s[1:] }
