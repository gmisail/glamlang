package ast

type Type interface {
	Equals(Type) bool
	String() string
}

type Node interface {
	String() string
	GetLine() int
}

type Expression interface {
	Node
	GetType() Type
}

type Statement interface {
	Node
}

type NodeMetadata struct {
	Line int
	Type Type
}

func CreateMetadata(line int) NodeMetadata {
	return NodeMetadata{line, nil}
}
