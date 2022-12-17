package codegen

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/codegen/backend"
)

type Compiler struct {
	//module *ir.Module
}

func (compiler *Compiler) compileStatement(statement ast.Statement) error {
	switch target := statement.(type) {
	case *ast.VariableDeclaration:
		return compiler.compileVariableDeclaration(target)
	}

	return fmt.Errorf("Unknown statement")
}

func (compiler *Compiler) compileVariableDeclaration(decl *ast.VariableDeclaration) error {
	name := decl.Name

	fmt.Println(backend.EmitRecord("Compiler"))
	fmt.Println(backend.EmitVariableDeclaration("int", name, backend.EmitBinary("+", "5", "10")))

	return nil
}

func (compiler *Compiler) dump() string {
	return ""
}

func Compile(statements []ast.Statement) string {
	compiler := Compiler{}

	for _, statement := range statements {
		compiler.compileStatement(statement)
	}

	return compiler.dump()
}
