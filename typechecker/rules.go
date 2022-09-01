package typechecker

import (
	"github.com/gmisail/glamlang/lexer"
)

var unaryRules = map[lexer.TokenType]map[string]bool{
	lexer.SUB: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.BANG: {
		"int":    false,
		"float":  false,
		"bool":   true,
		"string": false,
	},
}

var binaryRules = map[lexer.TokenType]map[string]bool{
	// math operations
	lexer.ADD: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": true,
	},
	lexer.SUB: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.MULT: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.DIV: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},

	// comparisons
	lexer.GT: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.GT_EQ: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.LT: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
	lexer.LT_EQ: {
		"int":    true,
		"float":  true,
		"bool":   false,
		"string": false,
	},
}

func HasBinaryRule(operation lexer.TokenType, variableType string) bool {
	if operation == lexer.EQUALITY || operation == lexer.NOT_EQUAL {
		return true
	}

	if op, opOk := binaryRules[operation]; opOk {
		if rule, ruleOk := op[variableType]; ruleOk {
			return rule
		}
	}

	return false
}

func HasUnaryRule(operation lexer.TokenType, variableType string) bool {
	if op, opOk := unaryRules[operation]; opOk {
		if rule, ruleOk := op[variableType]; ruleOk {
			return rule
		}
	}

	return false
}
