package ast

type Node interface {
	String() string
}

type Expression interface {
	Node
}

type Statement interface {
	Node
}

type TypeDefinition interface {
	Node
}
