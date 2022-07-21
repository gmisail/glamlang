package typechecker

import "github.com/gmisail/glamlang/ast"

func Check(statement ast.Statement) bool {
	return false
}

func TypeOf(expression *ast.Expression) {
	
}

func checkVariableDeclaration(v *ast.VariableDeclaration) bool {
	variableType := v.Type

	if v.Value != nil {
		

		return true
	}


	// shouldnt get here
	return true
}
