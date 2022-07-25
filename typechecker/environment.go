package typechecker

type Environment struct {
	Parent *Environment
	Values map[string]*Type
}

func CreateEnvironment(parent *Environment) *Environment {
	return &Environment{Parent: parent, Values: make(map[string]*Type)}
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
