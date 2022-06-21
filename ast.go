package main

type Node interface {
	String()
}

type Expression interface {
	Node
}

type Binary struct {
	Expression
	Left     Expression
	Right    Expression
	Operator TokenType
}

type Group struct {
	Expression
	Value Expression
}

type Literal struct {
	Expression
	Value interface{}
}
