package template

import (
	"go/types"
	"opt/internal/registry"
)

type MockData struct {
	InterfaceName string
	MockName      string
	TypeParams    []TypeParamData
	Methods       []MethodData
}

// MethodData is the data which represents a method on some interface.
type MethodData struct {
	Name    string
	Params  []ParamData
	Returns []ParamData
}
type TypeParamData struct {
	ParamData
	Constraint types.Type
}

// ParamData is the data which represents a parameter to some method of
// an interface.
type ParamData struct {
	Var      *registry.Var
	Variadic bool
}
