// Package svb defines interfaces for parsing
// and creating Simple Virtual Binary format files.
package svb

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
	Size     int
}

// Subroutine represents a subroutine defined in assembly.
type Subroutine struct {
	Name         string
	Address      uint16
	Instructions []Instruction
	Size         int
}

// SVB represents a Simple Virtual Binary format file.
type SVB struct {
	Constants   []Constant
	Subroutines []Subroutine
	MainAddress  uint16
}
