package vm

import (
	"writing-in-interpreter-in-go/src/monkey/code"
	"writing-in-interpreter-in-go/src/monkey/object"
)

type Frame struct {
	closure     *object.Closure
	ip          int
	basePointer int
}

func NewFrame(closure *object.Closure, basePointer int) *Frame {
	return &Frame{
		closure:     closure,
		ip:          -1,
		basePointer: basePointer,
	}
}
func (frame *Frame) Instructions() code.Instructions {
	return frame.closure.Function.Instructions
}
