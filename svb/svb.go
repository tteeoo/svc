// Package svb defines interfaces for parsing
//   and creating Simple Virtual Binary format files.
package svb

import (
	"github.com/tteeoo/svc/dat"
)

// Constant represents a constant defined in assembly.
type Constant struct {
	Name    string
	Address uint16
	Value   uint16
}

// Instruction represents an instruction defined in assembly.
type Instruction struct {
	Name     string
	Opcode   uint16
	Operands []uint16
}

// Size calculates the size of an Instruction.
func (i Instruction) Size() int {
	return dat.OpNameToSize[i.Name] + 1
}

// Subroutine represents a subroutine defined in assembly.
type Subroutine struct {
	Name         string
	Address      uint16
	Instructions []Instruction
}

// Size calculates the size of an Subroutine.
func (s Subroutine) Size() int {
	size := 0
	for _, i := range s.Instructions {
		size += i.Size()
	}
	return size
}

// SVB represents a Simple Virtual Binary formatted file.
type SVB struct {
	Constants   []Constant
	Subroutines []Subroutine
	MainAddress uint16
}

// Size calculates size of an SVB.
func (s SVB) Size() int {
	size := 0
	for _, sub := range s.Subroutines {
		size += sub.Size()
	}
	return len(s.Constants) + size
}
