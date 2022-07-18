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

func (v *VariableType) String() string {
	return fmt.Sprintf("(VariableType base: %s)", v.Base)
}

type FunctionType struct {
	TypeDefinition
	ArgumentTypes []TypeDefinition
	ReturnType    TypeDefinition
}

func (f *FunctionType) String() string {
	return fmt.Sprintf("(FunctionType )")
}

type Logical struct {
	Expression
	Left     Expression
	Right    Expression
	Operator TokenType
}

func (l *Logical) String() string {
	return fmt.Sprintf("(Logical)")
}

type VariableDeclaration struct {
	Statement
	Name  string
	Type  TypeDefinition
	Value Expression
}

func (v *VariableDeclaration) String() string {
	return fmt.Sprintf("(VariableDeclaration name: %s, type: %s, value: %s)", v.Name, v.Type.String(), v.Value.String())
}

type StructDeclaration struct {
	Statement
	Name      string
	Variables []VariableDeclaration
}

func (s *StructDeclaration) String() string {
	return fmt.Sprintf("(StructDeclaration body: )")
}

type ExpressionStatement struct {
	Statement
	Value Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("(ExpressionStatement body: )")
}

type BlockStatement struct {
	Statement
	Statements []Statement
}

func (b *BlockStatement) String() string {
	return fmt.Sprintf("(BlockStatement body: )")
}

type IfStatement struct {
	Statement
	Condition Expression
	Body      Statement
	ElseBody  Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("(IfStatement condition: %s, body: %s)", i.Condition.String(), i.Body.String())
}

type WhileStatement struct {
	Statement
	Condition Expression
	Body      Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("(WhileStatement condition: %s, body: %s)", w.Condition.String(), w.Body.String())
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

func (f *FunctionExpression) String() string {
	return fmt.Sprintf("(FunctionExpression body: %s)", f.Body.String())
}

type VariableExpression struct {
	Expression
	Value string
}

func (v *VariableExpression) String() string {
	return fmt.Sprintf("(VariableExpression %s)", v.Value)
}

type Literal struct {
	Expression
	Value interface{}
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal %s)", l.Value)
}
