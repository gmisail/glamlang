// Compiles an AST into LLVM IR.

package codegen

import (
	"github.com/gmisail/glamlang/ast"
	"tinygo.org/x/go-llvm"
)

type Compiler struct {
	module  llvm.Module
	context llvm.Context
	builder llvm.Builder
}

func (compiler *Compiler) compileStatement(statement ast.Statement) llvm.Value {
	switch target := statement.(type) {
	case *ast.ExpressionStatement:
		return compiler.compileExpression(target.Value)
	}

	return llvm.ConstNull(llvm.VoidType())
}

func (compiler *Compiler) compileExpression(expr ast.Expression) llvm.Value {
	switch target := expr.(type) {
	case *ast.Literal:
		return llvm.ConstFloat(llvm.FloatType(), 10.0)
	case *ast.Binary:
		return compiler.compileBinary(*target)
	}

	return llvm.ConstNull(llvm.VoidType())
}

func (compiler *Compiler) compileBinary(bin ast.Binary) llvm.Value {
	left := compiler.compileExpression(bin.Left)
	right := compiler.compileExpression(bin.Right)

	return compiler.builder.CreateFAdd(left, right, "add")
}

func Compile(statements []ast.Statement) {
	compiler := Compiler{
		llvm.NewModule("glam"),
		llvm.GlobalContext(),
		llvm.NewBuilder(),
	}

	for _, statement := range statements {
		compiler.compileStatement(statement).Dump()
	}
}
