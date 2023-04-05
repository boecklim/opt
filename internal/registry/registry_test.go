package registry_test

import (
	"fmt"
	"opt/internal/registry"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupStruct(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		registry, err := registry.New("testpackages/example/", "example")
		require.NoError(t, err)

		str, _, err := registry.LookupStruct("ExampleStruct")
		require.NoError(t, err)

		fmt.Printf("\n%s", str.Field(0).Name())
		fmt.Printf("\n%s", str.Field(1).Name())

	})
}
