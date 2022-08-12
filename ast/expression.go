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

type Unary struct {
	Expression
	NodeMetadata
	Value    Expression
	Operator lexer.TokenType
}

func (u *Unary) String() string {
	return fmt.Sprintf("(Unary %s %s)", lexer.TokenTypeToString(u.Operator), u.Value.String())
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

type Group struct {
	Expression
	NodeMetadata
	Value Expression
}

func (g *Group) String() string {
	return fmt.Sprintf("(Group %s)", g.Value.String())
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

type VariableExpression struct {
	Expression
	NodeMetadata
	Value string
}

func (v *VariableExpression) String() string {
	return fmt.Sprintf("(VariableExpression %s)", v.Value)
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

type GetExpression struct {
	Expression
	NodeMetadata
	Name   string
	Parent Expression
}

func (g *GetExpression) String() string {
	return fmt.Sprintf("(Get name: %s, parent: %s)", g.Name, g.Parent.String())
}

type Literal struct {
	Expression
	NodeMetadata
	Value interface{}
	Type  lexer.TokenType
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal %s)", l.Value)
}
