package vm

import (
	"writing-in-interpreter-in-go/src/monkey/code"
	"writing-in-interpreter-in-go/src/monkey/object"
)

type Frame struct {
	function    *object.CompiledFunction
	ip          int
	basePointer int
}

func NewFrame(function *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		function:    function,
		ip:          -1,
		basePointer: basePointer,
	}
}
func (frame *Frame) Instructions() code.Instructions {
	return frame.function.Instructions
}
