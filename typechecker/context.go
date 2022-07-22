package typechecker

type Context struct {
	Parent *Context
	Values map[string]*Type
}

func CreateContext(parent *Context) *Context {
	return &Context{Parent: parent, Values: make(map[string]*Type)}
}

func CreateGlobal() *Context {
	return CreateContext(nil)
}

/*
	Looks up the type of a variable if it exists.
*/
func (c *Context) Find(name string) (bool, *Type) {
	// check the current scope
	if variableType, ok := c.Values[name]; ok {
		return true, variableType
	}

	// if there are no more parents, then the variable
	// is not in scope / does not exist.
	if c.Parent == nil {
		return false, nil
	}

	// check in the upper scope
	return c.Parent.Find(name)
}

/*
	Adds a variable to the context. Returns false
	if the variable already exists in the current scope.
*/
func (c *Context) Add(variableName string, variableType *Type) bool {
	exists, _ := c.Find(variableName)

	if exists {
		return false
	}

	c.Values[variableName] = variableType

	return true
}
