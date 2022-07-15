package main

import "fmt"

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

type VariableType struct {
	TypeDefinition
	Base     string
	Optional bool
	SubType  *VariableType
}

type FunctionType struct {
	TypeDefinition
	ArgumentTypes []TypeDefinition
	ReturnType    TypeDefinition
}

type Logical struct {
	Expression
	Left     Expression
	Right    Expression
	Operator TokenType
}

type VariableDeclaration struct {
	Statement
	Name  string
	Type  TypeDefinition
	Value Expression
}

type ExpressionStatement struct {
	Statement
	Value Expression
}

type BlockStatement struct {
	Statement
	Statements []Statement
}

type IfStatement struct {
	Statement
	Condition Expression
	Body      Statement
	ElseBody  Statement
}

type WhileStatement struct {
	Statement
	Condition Expression
	Body      Statement
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

type FunctionExpression struct {
	Expression
	Parameters []string
	Body       Statement
}

type VariableExpression struct {
	Expression
	Value string
}

type Literal struct {
	Expression
	Value interface{}
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal %s)", l.Value)
}
