package ast

import (
	"fmt"
	"strings"
)

type VariableDeclaration struct {
	Statement
	Line  int
	Name  string
	Type  TypeDefinition
	Value Expression
}

func (v *VariableDeclaration) String() string {
	value := "nil"

	if v.Value != nil {
		value = v.Value.String()
	}

	return fmt.Sprintf("(VariableDeclaration name: %s, type: %s, value: %s)", v.Name, v.Type.String(), value)
}

type StructDeclaration struct {
	Statement
	Line      int
	Name      string
	Variables []VariableDeclaration
}

func (s *StructDeclaration) String() string {
	var builder strings.Builder

	builder.WriteString("(StructDeclaration [")

	for i, variable := range s.Variables {
		builder.WriteString(variable.String())

		if i != len(s.Variables)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("])")

	return builder.String()
}

type ExpressionStatement struct {
	Statement
	Line  int
	Value Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("(ExpressionStatement body: %s)", e.Value.String())
}

type BlockStatement struct {
	Statement
	Line       int
	Statements []Statement
}

func (b *BlockStatement) String() string {
	return fmt.Sprintf("(BlockStatement body: )")
}

type IfStatement struct {
	Statement
	Line      int
	Condition Expression
	Body      Statement
	ElseBody  Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("(IfStatement condition: %s, body: %s)", i.Condition.String(), i.Body.String())
}

type WhileStatement struct {
	Statement
	Line      int
	Condition Expression
	Body      Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("(WhileStatement condition: %s, body: %s)", w.Condition.String(), w.Body.String())
}

type ReturnStatement struct {
	Statement
	Line  int
	Value Expression
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("(ReturnStatement value: %s)", r.Value.String())
}
