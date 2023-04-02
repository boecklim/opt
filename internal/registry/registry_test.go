package registry_test

import (
	"opt/internal/registry"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupStruct(t *testing.T) {
	t.Run("test", func(t *testing.T) {

		registry, err := registry.New("./", "main")
		require.NoError(t, err)

		_, _, err = registry.LookupInterface("ExampleInterface")
		require.NoError(t, err)

		_, _, err = registry.LookupStruct("ExampleStruct")
		require.NoError(t, err)

	})
}
