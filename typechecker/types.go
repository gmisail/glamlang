package typechecker

type Type interface { }

type TypeInt struct {
	Type
}

type TypeFloat struct {
	Type
}

type TypeString struct {
	Type
}

type TypeBool struct {
	Type
}

type TypeVariable struct {
	Type
	Value Type
}