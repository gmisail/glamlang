package main

import "fmt"

type Node interface {
	String() string
}

type Expression interface {
	Node
}

type Unary struct {
	Expression
	Value    Expression
	Operator TokenType
}

func (u *Unary) String() string {
	return fmt.Sprintf("(Unary %s %s)", tokenTypeToString(u.Operator), u.Value.String())
}

type Binary struct {
	Expression
	Left     Expression
	Right    Expression
	Operator TokenType
}

func (b *Binary) String() string {
	return fmt.Sprintf("(Binary %s %s)", b.Left.String(), b.Right.String())
}

type Group struct {
	Expression
	Value Expression
}

func (g *Group) String() string {
	return fmt.Sprintf("(Group %s)", g.Value.String())
}

type Literal struct {
	Expression
	Value interface{}
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal %s)", l.Value)
}
