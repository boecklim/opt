package opt

import (
	"bytes"
	"errors"
	"io"
	"opt/internal/registry"

	"opt/internal/template"
)

type Mocker struct {
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

// New makes a new Mocker for the specified package directory.
func New(cfg Config) (*Mocker, error) {
	reg, err := registry.New(cfg.SrcDir, cfg.PkgName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	return &Mocker{
		cfg:      cfg,
		registry: reg,
		tmpl:     tmpl,
	}, nil
}

// Mock generates a mock for the specified interface name.
func (m *Mocker) Mock(w io.Writer, namePairs ...string) error {
	if len(namePairs) == 0 {
		return errors.New("must specify one struct")
	}

	data := template.Data{
		PkgName: m.mockPkgName(),
	}

	var buf bytes.Buffer
	if err := m.tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return nil
}

func (m *Mocker) mockPkgName() string {
	if m.cfg.PkgName != "" {
		return m.cfg.PkgName
	}

	return m.registry.SrcPkgName()
}
