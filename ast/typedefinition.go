package ast

import (
	"fmt"
	"strings"

	"github.com/gmisail/glamlang/lexer"
)

type VariableType struct {
	Type
	Base     string
	Optional bool
	SubType  *VariableType
}

func (v *VariableType) String() string {
	return v.Base
}

type FunctionType struct {
	Type
	Parameters []Type
	ReturnType Type
}

var internalTypes = map[string]Type{
	"int":    &VariableType{Base: "int", Optional: false},
	"float":  &VariableType{Base: "float", Optional: false},
	"string": &VariableType{Base: "string", Optional: false},
	"bool":   &VariableType{Base: "bool", Optional: false},
	"null":   &VariableType{Base: "null", Optional: false},
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
		if v.Base != target.Base || v.Optional != target.Optional {
			return false
		}
	case *FunctionType:
		return false
	}

	return true
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
	return &VariableType{Base: name, Optional: isOptional}
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

func CreateTypeFrom(t Type) Type {
	switch targetType := t.(type) {
	case *VariableType:
		return CreateVariableType(targetType.Base, targetType.Optional)
	case *FunctionType:
		parameters := make([]Type, len(targetType.Parameters))

		for i, param := range targetType.Parameters {
			parameters[i] = CreateTypeFrom(param)
		}

		return CreateFunctionType(parameters, CreateTypeFrom(targetType.ReturnType))
	}

	return nil
}

func IsInternalType(target Type) (bool, Type) {
	if targetType, ok := target.(*VariableType); ok {
		//isOptional := targetType.Optional

		if internalType, typeExists := internalTypes[targetType.Base]; typeExists {
			return true, internalType
		}
	}

	return false, nil
}
