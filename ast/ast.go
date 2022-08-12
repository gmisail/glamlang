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

type Type interface {
	Node
	Equals(Type) bool
}

type NodeMetadata struct {
	Line int
	Type Type
}

func CreateMetadata(line int) NodeMetadata {
	return NodeMetadata{line, nil}
}
