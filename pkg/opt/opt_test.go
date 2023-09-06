package opt

import (
	"bytes"
	"fmt"
	"opt/pkg/opt/testpackages/example"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOpt(t *testing.T) {
	m, err := New(Config{SrcDir: "testpackages/example/"})
	require.NoError(t, err)

	var buf bytes.Buffer
	err = m.Generate(&buf, "ExampleStruct")
	require.NoError(t, err)

	// s := buf.String()
	// fmt.Println(s)

	err = os.WriteFile("testpackages/example/constructor.go", buf.Bytes(), 0o600)
	require.NoError(t, err)
}

func TestConstructor(t *testing.T) {
	testString := "hello"
	example := example.New(
		"firstMember",
		0,
		time.Now(),
		example.WithAPointer(&testString),
		example.WithASlice([]float64{0.0, 1.1}),
	)

	fmt.Println(example)
}
