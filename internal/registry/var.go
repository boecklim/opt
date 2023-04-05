package registry

import (
	"go/types"
	"strings"
)

// varNameForType generates a name for the variable using the type
// information.
//
// Examples:
// - string -> s
// - int -> n
// - chan int -> intCh
// - []a.MyType -> myTypes
// - map[string]int -> stringToInt
// - error -> err
// - a.MyType -> myType
func varNameForType(t types.Type) string {
	nestedType := func(t types.Type) string {
		if t, ok := t.(*types.Basic); ok {
			return deCapitalise(t.String())
		}
		return varNameForType(t)
	}

	switch t := t.(type) {
	case *types.Named:
		if t.Obj().Name() == "error" {
			return "err"
		}

		name := deCapitalise(t.Obj().Name())
		if name == t.Obj().Name() {
			name += "MoqParam"
		}

		return name

	case *types.Basic:
		return basicTypeVarName(t)

	case *types.Array:
		return nestedType(t.Elem()) + "s"

	case *types.Slice:
		return nestedType(t.Elem()) + "s"

	case *types.Struct: // anonymous struct
		return "val"

	case *types.Pointer:
		return varNameForType(t.Elem())

	case *types.Signature:
		return "fn"

	case *types.Interface: // anonymous interface
		return "ifaceVal"

	case *types.Map:
		return nestedType(t.Key()) + "To" + capitalise(nestedType(t.Elem()))

	case *types.Chan:
		return nestedType(t.Elem()) + "Ch"
	}

	return "v"
}

func basicTypeVarName(b *types.Basic) string {
	switch b.Info() {
	case types.IsBoolean:
		return "b"

	case types.IsInteger:
		return "n"

	case types.IsFloat:
		return "f"

	case types.IsString:
		return "s"
	}

	return "v"
}

func capitalise(s string) string   { return strings.ToUpper(s[:1]) + s[1:] }
func deCapitalise(s string) string { return strings.ToLower(s[:1]) + s[1:] }
