package object

import "bytelang/src/bytelang/code"

type CompiledFunction struct {
	Instructions       code.Instructions
	ParameterArity     int
	LocalVariableArity int
}
