// Creates basic C elements like variable declarations & struct definitions.

package backend

import "fmt"

func EmitRecord(name string) string {
	return fmt.Sprintf("typedef struct { } %s;", name)
}

func EmitVariableDeclaration(variableType string, name string, value string) string {
	return fmt.Sprintf("%s %s = %s;", variableType, name, value)
}

func EmitFunctionDeclaration(variableType string, name string, body string) string {
	return fmt.Sprintf("%s %s(){ %s }", variableType, name, body)
}

func EmitUnary(symbol string, value string) string {
	return symbol + value
}

func EmitBinary(symbol string, left string, right string) string {
	return left + symbol + right
}
