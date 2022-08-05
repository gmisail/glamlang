package typechecker

import (
	"fmt"
	"strings"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

type Type interface {
	Equals(Type) bool
	String() string
}

type VariableType struct {
	Type
	Name       string
	IsOptional bool
}

type FunctionType struct {
	Type
	Parameters []Type
	ReturnType Type
}

var internalTypes = map[string]Type{
	"int":    &VariableType{Name: "int", IsOptional: false},
	"float":  &VariableType{Name: "float", IsOptional: false},
	"string": &VariableType{Name: "string", IsOptional: false},
	"bool":   &VariableType{Name: "bool", IsOptional: false},
	"null":   &VariableType{Name: "null", IsOptional: false},
}

func (v *VariableType) Equals(otherType Type) bool {
	if otherType == nil {
		return false
	}

	/*
		int == int 	 	   (yes)
		bool == bool 	   (yes)
		int == int?  	   (no)
		string? == string? (yes)
		(int, int) -> int == (int, float) -> int (no)
		(int, int) -> int == (int, int) -> int   (yes)
	*/

	switch target := otherType.(type) {
	case *VariableType:
		if v.Name != target.Name || v.IsOptional != target.IsOptional {
			return false
		}
	case *FunctionType:
		return false
	}

	return true
}

func (v *VariableType) String() string {
	return v.Name
}

func (f *FunctionType) Equals(otherType Type) bool {
	if otherType == nil {
		return false
	}

	/*
		If the otherType is a variable type (i.e. int, string, etc...)
		automatically reject it since a function signature =/= variable type.
	*/
	switch target := otherType.(type) {
	case *VariableType:
		return false
	case *FunctionType:
		// validate length of parameters
		if len(target.Parameters) != len(f.Parameters) {
			return false
		}

		// validate that every parameter matches
		for i, param := range f.Parameters {
			if !param.Equals(target.Parameters[i]) {
				return false
			}
		}

		if !target.ReturnType.Equals(f.ReturnType) {
			return false
		}
	}

	return true
}

func (f *FunctionType) String() string {
	var builder strings.Builder

	for i, param := range f.Parameters {
		builder.WriteString(param.String())

		if i != len(f.Parameters)-1 {
			builder.WriteString(", ")
		}
	}

	return fmt.Sprintf("(%s) -> %s", builder.String(), f.ReturnType.String())
}

func CreateVariableType(name string, isOptional bool) *VariableType {
	return &VariableType{Name: name, IsOptional: isOptional}
}

func CreateFunctionType(parameters []Type, returnType Type) *FunctionType {
	return &FunctionType{Parameters: parameters, ReturnType: returnType}
}

func CreateTypeFromLiteral(literalType lexer.TokenType) Type {
	typeName := strings.ToLower(lexer.TokenTypeToString(literalType))

	if internalType, typeExists := internalTypes[typeName]; typeExists {
		return internalType
	}

	return internalTypes["null"]
}

func CreateTypeFrom(typeDefinition ast.TypeDefinition) Type {
	switch targetType := typeDefinition.(type) {
	case *ast.VariableType:
		return CreateVariableType(targetType.Base, targetType.Optional)
	case *ast.FunctionType:
		parameters := make([]Type, len(targetType.ArgumentTypes))

		for i, param := range targetType.ArgumentTypes {
			parameters[i] = CreateTypeFrom(param)
		}

		return CreateFunctionType(parameters, CreateTypeFrom(targetType.ReturnType))
	}

	return nil
}

func IsInternalType(target ast.TypeDefinition) (Type, error) {
	if targetType, ok := target.(*ast.VariableType); ok {
		isOptional := targetType.Optional

		if internalType, typeExists := internalTypes[targetType.Base]; typeExists {
			return isOptional, internalType
		}
	}

	return false, nil
}
