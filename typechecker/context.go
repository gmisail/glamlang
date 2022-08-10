package typechecker

import (
	"fmt"

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
func (c *Context) Find(name string) (bool, *ast.Type) {
	isValid, variableType := c.environment.Find(name)

	if isValid {
		return isValid, variableType
	}

	return false, nil
}

/*
	Adds variable to the current scope.
*/
func (c *Context) Add(variableName string, variableType *ast.Type) bool {
	return c.environment.Add(variableName, variableType)
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

	for k := range c.environment.Values {
		fmt.Println("pop", k)
	}

	fmt.Println("pop scope")

	c.environment = c.environment.Parent
}
