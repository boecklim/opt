package opt

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOpt(t *testing.T) {
	m, err := New(Config{SrcDir: "testpackages/example", PkgName: "example"})
	if err != nil {
		t.Fatalf("opt.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "ExampleStruct")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	fmt.Println(s)
}
