package context

import (
	"github.com/gmisail/glamlang/ast"
)

type Environment struct {
	Parent *Environment
	Values map[string]*ast.Type
	Types  map[string]ast.RecordType
}

func CreateEnvironment(parent *Environment) *Environment {
	return &Environment{
		Parent: parent,
		Values: make(map[string]*ast.Type),
		Types:  make(map[string]ast.RecordType),
	}
}

/*
Looks up the type of a variable if it exists.
*/
func (e *Environment) FindVariable(name string) (bool, *ast.Type) {
	// check the current scope
	if variableType, ok := e.Values[name]; ok {
		return true, variableType
	}

	if e.Parent == nil {
		return false, nil
	}

	return e.Parent.FindVariable(name)
}

/*
Adds a variable to the context. Returns false
if the variable already exists in the current scope.
*/
func (e *Environment) AddVariable(variableName string, variableType *ast.Type) bool {
	exists, _ := e.FindVariable(variableName)

	if exists {
		return false
	}

	e.Values[variableName] = variableType

	return true
}

/*
Returns if a custom type exists in the current context.
*/
func (e *Environment) TypeExists(typeName string) bool {
	exists, _ := e.FindType(typeName)
	return exists
}

/*
Adds a custom type (i.e. struct) if it does not exist already. Returns
true if it was added as well as an environment which represents the
variables within the struct.
*/
func (e *Environment) AddType(typeName string, record ast.RecordType) bool {
	if e.TypeExists(typeName) {
		return false
	}

	e.Types[typeName] = record

	return true
}

func (e *Environment) FindType(typeName string) (bool, *ast.RecordType) {
	if customType, ok := e.Types[typeName]; ok {
		return true, &customType
	}

	if e.Parent == nil {
		return false, nil
	}

	return e.Parent.FindType(typeName)
}
