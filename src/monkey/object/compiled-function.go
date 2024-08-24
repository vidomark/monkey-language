package object

import "writing-in-interpreter-in-go/src/monkey/code"

type CompiledFunction struct {
	Instructions       code.Instructions
	ParameterArity     int
	LocalVariableArity int
}
