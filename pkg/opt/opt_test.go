package opt

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpt(t *testing.T) {
	m, err := New(Config{SrcDir: "testpackages/example/", PkgName: "example"})
	require.NoError(t, err)

	var buf bytes.Buffer
	err = m.Generate(&buf, "ExampleStruct")
	require.NoError(t, err)

	// s := buf.String()
	// fmt.Println(s)

	err = os.WriteFile("testpackages/example/example_out.go", buf.Bytes(), 0o600)
	require.NoError(t, err)
}
