package vm

import (
	"bytelang/src/bytelang/code"
	"bytelang/src/bytelang/object"
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
