package ast

import "tinygo.org/x/go-llvm"

func (l *Logical) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (u *Unary) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (b *Binary) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())
}

func (g *Group) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (f *FunctionExpression) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (v *VariableExpression) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (f *FunctionCall) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (g *GetExpression) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (l *Literal) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (v *VariableDeclaration) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (s *RecordDeclaration) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (e *ExpressionStatement) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (b *BlockStatement) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (i *IfStatement) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (w *WhileStatement) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())

}

func (r *ReturnStatement) Codegen() llvm.Value {
	return llvm.ConstNull(llvm.VoidType())
}
