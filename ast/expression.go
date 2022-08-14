package ast

import (
	"fmt"
	"strings"

	"github.com/gmisail/glamlang/lexer"
)

type Logical struct {
	Expression
	NodeMetadata
	Left     Expression
	Right    Expression
	Operator lexer.TokenType
}

func (l *Logical) String() string {
	return fmt.Sprintf(
		"(Logical op: %s, left: %s, right: %s)",
		lexer.TokenTypeToString(l.Operator),
		l.Left.String(),
		l.Right.String(),
	)
}

func (l *Logical) GetLine() int {
	return l.NodeMetadata.Line
}

func (l *Logical) GetType() Type {
	return l.NodeMetadata.Type
}

type Unary struct {
	Expression
	NodeMetadata
	Value    Expression
	Operator lexer.TokenType
}

func (u *Unary) String() string {
	return fmt.Sprintf("(Unary %s %s)", lexer.TokenTypeToString(u.Operator), u.Value.String())
}

func (u *Unary) GetLine() int {
	return u.NodeMetadata.Line
}

func (u *Unary) GetType() Type {
	return u.NodeMetadata.Type
}

type Binary struct {
	Expression
	NodeMetadata
	Left     Expression
	Right    Expression
	Operator lexer.TokenType
}

func (b *Binary) String() string {
	return fmt.Sprintf("(Binary %s %s)", b.Left.String(), b.Right.String())
}

func (b *Binary) GetLine() int {
	return b.NodeMetadata.Line
}

func (b *Binary) GetType() Type {
	return b.NodeMetadata.Type
}

type Group struct {
	Expression
	NodeMetadata
	Value Expression
}

func (g *Group) String() string {
	return fmt.Sprintf("(Group %s)", g.Value.String())
}

func (g *Group) GetLine() int {
	return g.NodeMetadata.Line
}

func (g *Group) GetType() Type {
	return g.NodeMetadata.Type
}

type FunctionExpression struct {
	Expression
	NodeMetadata
	Parameters []VariableDeclaration
	Body       Statement
	ReturnType Type
}

func (f *FunctionExpression) String() string {
	return fmt.Sprintf("(FunctionExpression body: %s)", f.Body.String())
}

func (f *FunctionExpression) GetLine() int {
	return f.NodeMetadata.Line
}

func (f *FunctionExpression) GetType() Type {
	return f.NodeMetadata.Type
}

type VariableExpression struct {
	Expression
	NodeMetadata
	Value string
}

func (v *VariableExpression) String() string {
	return fmt.Sprintf("(VariableExpression %s)", v.Value)
}

func (v *VariableExpression) GetLine() int {
	return v.NodeMetadata.Line
}

func (v *VariableExpression) GetType() Type {
	return v.NodeMetadata.Type
}

type FunctionCall struct {
	Expression
	NodeMetadata
	Callee    Expression
	Arguments []Expression
}

func (f *FunctionCall) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("(FunctionCall callee: %s, arguments: [", f.Callee.String()))

	for i, argument := range f.Arguments {
		builder.WriteString(argument.String())

		if i != len(f.Arguments)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("])")

	return builder.String()
}

func (f *FunctionCall) GetLine() int {
	return f.NodeMetadata.Line
}

func (f *FunctionCall) GetType() Type {
	return f.NodeMetadata.Type
}

type GetExpression struct {
	Expression
	NodeMetadata
	Name   string
	Parent Expression
}

func (g *GetExpression) String() string {
	return fmt.Sprintf("(Get name: %s, parent: %s)", g.Name, g.Parent.String())
}

func (g *GetExpression) GetLine() int {
	return g.NodeMetadata.Line
}

func (g *GetExpression) GetType() Type {
	return g.NodeMetadata.Type
}

type Literal struct {
	Expression
	NodeMetadata
	Value       interface{}
	LiteralType lexer.TokenType
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal %s)", l.Value)
}

func (l *Literal) GetLine() int {
	return l.NodeMetadata.Line
}

func (l *Literal) GetType() Type {
	return l.NodeMetadata.Type
}
