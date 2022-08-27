package ast

import (
	"fmt"
	"strings"
)

type VariableDeclaration struct {
	Statement
	NodeMetadata
	Name  string
	Type  Type
	Value Expression
}

func (v *VariableDeclaration) String() string {
	value := "nil"

	if v.Value != nil {
		value = v.Value.String()
	}

	return fmt.Sprintf(
		"(VariableDeclaration name: %s, type: %s, value: %s)",
		v.Name,
		v.Type.String(),
		value,
	)
}

type RecordDeclaration struct {
	Statement
	NodeMetadata
	Name     string
	Record   RecordType
	Inherits string
}

func (s *RecordDeclaration) String() string {
	var builder strings.Builder

	builder.WriteString("(RecordDeclaration ")
	builder.WriteString(s.Record.String())
	builder.WriteString(")")

	return builder.String()
}

type ExpressionStatement struct {
	Statement
	NodeMetadata
	Value Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("(ExpressionStatement body: %s)", e.Value.String())
}

type BlockStatement struct {
	Statement
	NodeMetadata
	Statements []Statement
}

func (b *BlockStatement) String() string {
	return "(BlockStatement body: )"
}

type IfStatement struct {
	Statement
	NodeMetadata
	Condition Expression
	Body      Statement
	ElseBody  Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf(
		"(IfStatement condition: %s, body: %s)",
		i.Condition.String(),
		i.Body.String(),
	)
}

type WhileStatement struct {
	Statement
	NodeMetadata
	Condition Expression
	Body      Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf(
		"(WhileStatement condition: %s, body: %s)",
		w.Condition.String(),
		w.Body.String(),
	)
}

type ReturnStatement struct {
	Statement
	NodeMetadata
	Value Expression
	Type  interface{}
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("(ReturnStatement value: %s)", r.Value.String())
}
