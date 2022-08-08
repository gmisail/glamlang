package ast

import "fmt"

type VariableType struct {
	TypeDefinition
	Base     string
	Optional bool
	SubType  *VariableType
}

func (v *VariableType) String() string {
	return fmt.Sprintf("(VariableType base: %s)", v.Base)
}

type FunctionType struct {
	TypeDefinition
	ArgumentTypes []TypeDefinition
	ReturnType    TypeDefinition
}

func (f *FunctionType) String() string {
	return fmt.Sprintf("(FunctionType )")
}
