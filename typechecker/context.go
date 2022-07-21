package typechecker

type Context {
	parent *Context
	values map[string]interface{}
}

func CreateContext(parent *Context) *Context {
	return &Context{ parent: parent, values: make(map[string]interface{}) } 
}

func CreateGlobal() *Context {
	return CreateEnvironment(nil)
}

/*
	Looks up the type of a variable if it exists.
*/
func (c *Context) Find(name string) (bool, Type) {
	if variableType, ok := c.values[name]; ok {
		return true, variableType
	}

	return false, nil
}