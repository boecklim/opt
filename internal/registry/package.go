package registry

import "go/types"

type Package struct {
	pkg *types.Package

	Alias string
}
