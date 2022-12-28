// Compiles an AST into LLVM IR.

package codegen

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
)

type Compiler struct {
	context Context
	emitter Emitter
}

func Compile(emitter Emitter, statements []ast.Statement) {
	compiler := Compiler{
		Context{},
		emitter,
	}

	// fmt.Println(emitter.EmitUnary("-", "7"))
	// fmt.Println(emitter.EmitBinary("*", "10", "100"))

	// fmt.Println(emitter.EmitVariableDeclaration("temp_value", "int", emitter.EmitBinary("*", emitter.EmitUnary("-", "7"), "100")))

	// fmt.Println(emitter.EmitRecordDeclaration("Token", map[string]string{
	// 	"literal":  "string",
	// 	"location": "int",
	// }))

	// fmt.Println(emitter.EmitVariableDeclaration("x", "int", emitter.EmitBinary("*", "10", emitter.EmitUnary("-", "50"))))

	for _, statement := range statements {
		fmt.Println(compiler.CompileStatement(statement))
	}
}
