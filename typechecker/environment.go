package typechecker

type Environment struct {
	Parent *Environment
	Values map[string]*Type
	Types  map[string]*Environment
}

func CreateEnvironment(parent *Environment) *Environment {
	return &Environment{Parent: parent, Values: make(map[string]*Type), Types: make(map[string]*Environment)}
}

/*
	Looks up the type of a variable if it exists.
*/
func (e *Environment) Find(name string) (bool, *Type) {
	// check the current scope
	if variableType, ok := e.Values[name]; ok {
		return true, variableType
	}

	if e.Parent == nil {
		return false, nil
	}

	return e.Parent.Find(name)
}

/*
	Adds a variable to the context. Returns false
	if the variable already exists in the current scope.
*/
func (e *Environment) Add(variableName string, variableType *Type) bool {
	exists, _ := e.Find(variableName)

	if exists {
		return false
	}

	e.Values[variableName] = variableType

	return true
}

/*
	Returns if a custom type exists in the current context.
*/
func (e *Environment) CustomTypeExists(typeName string) bool {
	if _, ok := e.Types[typeName]; ok {
		return true
	}

	return false
}

/*
	Adds a custom type (i.e. struct) if it does not exist already. Returns
	true if it was added as well as an environment which represents the
	variables within the struct.
*/
func (e *Environment) AddType(typeName string) (bool, *Environment) {
	if e.CustomTypeExists(typeName) {
		return false, nil
	}

	environment := CreateEnvironment(nil)
	e.Types[typeName] = environment

	return true, environment
}
