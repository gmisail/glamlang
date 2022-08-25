package context

import (
	"github.com/gmisail/glamlang/ast"
)

type Context struct {
	environment *Environment
}

func CreateContext() *Context {
	return &Context{environment: CreateEnvironment(nil)}
}

/*
Looks up the type of a variable if it exists.
*/
func (c *Context) FindVariable(name string) (bool, *ast.Type) {
	isValid, variableType := c.environment.FindVariable(name)

	if isValid {
		return isValid, variableType
	}

	return false, nil
}

/*
Adds variable to the current scope.
*/
func (c *Context) Add(variableName string, variableType *ast.Type) bool {
	return c.environment.AddVariable(variableName, variableType)
}

func (c *Context) AddType(typeName string, recordType ast.RecordType) bool {
	return c.environment.AddType(typeName, recordType)
}

func (c *Context) FindType(typeName string) (bool, *ast.RecordType) {
	return c.environment.FindType(typeName)
}

func (c *Context) TypeExists(typeName string) bool {
	return c.environment.TypeExists(typeName)
}

/*
Creates and enters a new environment.
*/
func (c *Context) EnterScope() {
	c.environment = CreateEnvironment(c.environment)
}

/*
Pops the current environment and replaces it
with the parent. If the parent environment does
no exist, i.e. we are already in the global scope,
then just ignore it.
*/
func (c *Context) ExitScope() {
	if c.environment.Parent == nil {
		return
	}

	c.environment = c.environment.Parent
}
