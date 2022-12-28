// Emits Javascript constructs such that we can construct Javascript code via Go.

package backend

import (
	"fmt"
	"strings"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

type NativeEmitter struct{}

func (e NativeEmitter) ResolveLiteral(node *ast.Literal) string {
	switch node.LiteralType {
	case lexer.BOOL:
		value := node.Value.(bool)

		if value {
			return "true"
		}

		return "false"
	case lexer.INT:
		return fmt.Sprintf("%s", node.Value.(string))
	case lexer.FLOAT:
		return fmt.Sprintf("%s", node.Value.(string))
	case lexer.STRING:
		return fmt.Sprintf("\"%s\"", node.Value.(string))
	}

	return ""
}

func (e NativeEmitter) EmitGroup(operation string) string {
	return "(" + operation + ")"
}

func (e NativeEmitter) EmitUnary(operation string, value string) string {
	return e.EmitGroup(operation + value)
}

func (e NativeEmitter) EmitBinary(operation string, left string, right string) string {
	return e.EmitGroup(left + operation + right)
}

func (e NativeEmitter) EmitVariableDeclaration(variableName string, variableType string, value string) string {
	return fmt.Sprintf("%s %s = %s;", variableType, variableName, value)
}

func (e NativeEmitter) EmitRecordDeclaration(recordName string, fields map[string]string) string {
	var builder strings.Builder
	builder.WriteString("typedef struct {\n")

	for field, fieldType := range fields {
		builder.WriteString(fmt.Sprintf("%s %s;\n", fieldType, field))
	}

	builder.WriteString(fmt.Sprintf("} %s;", recordName))

	return builder.String()
}
