package typechecker

import "github.com/gmisail/glamlang/ast"

type Type struct {
	Name       string
	IsOptional bool
}

func CreateType(name string, isOptional bool) *Type {
	return &Type{Name: name, IsOptional: isOptional}
}

func CreateTypeFrom(typeDefinition ast.TypeDefinition) *Type {
	switch targetType := typeDefinition.(type) {
	case *ast.VariableType:
		return CreateType(targetType.Base, targetType.Optional)
	case *ast.FunctionType:
		return nil // TODO: figure out how to represent function types
	}

	return nil
}

func (t *Type) Equals(otherType *Type) bool {
	/*
		int == int 	 	   (yes)
		bool == bool 	   (yes)
		int == int?  	   (no)
		string? == string? (yes)
		(int, int) -> int == (int, float) -> int (no)
		(int, int) -> int == (int, int) -> int   (yes)
	*/
	if t.Name != otherType.Name || t.IsOptional != otherType.IsOptional {
		return false
	}

	return true
}

func IsInternalType(target ast.TypeDefinition) (bool, *Type) {
	if targetType, ok := target.(*ast.VariableType); ok {
		isOptional := targetType.Optional

		switch targetType.Base {
		case "int":
			return true, CreateType("int", isOptional)
		case "float":
			return true, CreateType("float", isOptional)
		case "bool":
			return true, CreateType("bool", isOptional)
		default:
			return false, nil
		}
	}

	// functions?

	return false, nil
}
