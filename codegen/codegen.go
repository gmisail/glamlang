package codegen

import (
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (c *Compiler) CompileExpression(node ast.Expression) string {
	switch expr := node.(type) {
	case *ast.Literal:
		return c.compileLiteral(expr)
	case *ast.Group:
		return c.emitter.EmitGroup(c.CompileExpression(expr.Value))
	case *ast.Binary:
		return c.compileBinary(expr)
	case *ast.Unary:
		return c.emitter.EmitUnary("-", c.CompileExpression(expr.Value))
	default:
		return ""
	}

	return "nil"
}

func (c *Compiler) CompileStatement(node ast.Statement) string {
	switch statement := node.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(statement.Value)
	default:
		return "idk"
	}

	return "eof"
}

func (c *Compiler) compileLiteral(literal *ast.Literal) string {
	return c.emitter.ResolveLiteral(literal)
}

func (c *Compiler) compileBinary(binary *ast.Binary) string {
	var operator string

	switch binary.Operator {
	case lexer.ADD:
		operator = "+"
	case lexer.SUB:
		operator = "-"
	case lexer.MULT:
		operator = "*"
	case lexer.DIV:
		operator = "/"
	case lexer.LT:
		operator = "<"
	case lexer.LT_EQ:
		operator = "<="
	case lexer.GT:
		operator = ">"
	case lexer.GT_EQ:
		operator = ">="
	case lexer.EQUALITY:
		operator = "=="
	case lexer.NOT_EQUAL:
		operator = "!="
	default:
		panic("Unknown binary case.")
	}

	return c.emitter.EmitBinary(operator, c.CompileExpression(binary.Left), c.CompileExpression(binary.Right))
}

func (c *Compiler) compileLogical(logical *ast.Logical) string {
	var operator string

	switch logical.Operator {
	case lexer.LT:
		operator = "<"
	case lexer.LT_EQ:
		operator = "<="
	case lexer.GT:
		operator = ">"
	case lexer.GT_EQ:
		operator = ">="
	case lexer.EQUALITY:
		operator = "=="
	case lexer.NOT_EQUAL:
		operator = "!="
	default:
		panic("Unknown logical case.")
	}

	return c.emitter.EmitBinary(operator, c.CompileExpression(logical.Left), c.CompileExpression(logical.Right))
}
