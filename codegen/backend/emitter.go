// Emits Javascript constructs such that we can construct Javascript code via Go.

package backend

import (
	"fmt"
	"strings"
)

type NativeEmitter struct{}

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
