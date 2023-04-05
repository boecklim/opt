package opt

import (
	"bytes"
	"errors"
	"go/format"
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
	SrcDir  string
	PkgName string
	// Formatter  string
	// StubImpl   bool
	// SkipEnsure bool
	// WithResets bool
}

func New(cfg Config) (*Generator, error) {
	reg, err := registry.New(cfg.SrcDir, cfg.PkgName)
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

	for i := 0; i < structType.NumFields(); i++ {
		members[i] = template.Member{
			Name:       structType.Field(i).Name(),
			Type:       structType.Field(i).Type().String(),
			StructName: structName,
		}
	}

	data := template.Data{
		PkgName:    m.mockPkgName(),
		StructName: structName,
		Members:    members,
	}

	srcPkgName := m.registry.SrcPkgName()
	mockPkgName := m.mockPkgName()
	if srcPkgName != mockPkgName {
		data.SrcPkgQualifier = m.registry.SrcPkgName() + "."
		imprt := m.registry.AddImport(m.registry.SrcPkg())
		data.SrcPkgQualifier = imprt.Qualifier() + "."
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

func (m *Generator) mockPkgName() string {
	if m.cfg.PkgName != "" {
		return m.cfg.PkgName
	}

	return m.registry.SrcPkgName()
}
