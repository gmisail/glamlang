package codegen

import "github.com/gmisail/glamlang/codegen/backend"

type Emitter interface {
	EmitGroup(string) string
	EmitUnary(string, string) string
	EmitBinary(string, string, string) string
	EmitRecordDeclaration(string, map[string]string) string
	EmitVariableDeclaration(string, string, string) string
}

func GetNativeBackend() backend.NativeEmitter {
	return backend.NativeEmitter{}
}
