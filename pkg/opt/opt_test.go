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
	example := example.New(
		example.WithFirstMember("hello"),
		example.WithSecondMember(100),
		example.WithTimestamp(time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)),
	)

	fmt.Println(example)
}
