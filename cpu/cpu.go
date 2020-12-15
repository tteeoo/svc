// Package cpu implements structs and functions for
// representing and operating the virtual CPU.
package cpu

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
)

// CPU is a basic implementation of a CPU.
type CPU struct {
	// mem is the memory device used by the CPU.
	Mem *mem.RAM
	// regs maps numbers to regsiters.
	Regs map[uint16]*Register
	// ops maps opcode names to opcodes.
	Ops map[string]uint16
}

// NewCPU returns a pointer to a newly initialized CPU.
func NewCPU() *CPU {
	regs := make(map[uint16]*Register)
	// General purpose registers
	for _, i := range []uint16{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07} {
		regs[i] = NewRegister(i)
	}
	// Program counter
	regs[0x08] = NewRegister(0x08)
	// Accumulator
	regs[0x09] = NewRegister(0x09)
	return &CPU{
		Mem:  mem.NewRAM(),
		Regs: regs,
		Ops: map[string]uint16{
			"nop": 0x00,
			"cop": 0x01,
			"cpl": 0x02,
			"str": 0x03,
			"ldr": 0x04,
		},
	}
}

// GetMem returns the CPU's memory device.
func (c *CPU) GetMem() *mem.RAM {
	return c.Mem
}

// GetOp returns to opcode whose name is provided.
// Returns 0x00 (nop) if the name is not defined.
func (c *CPU) GetOp(name string) uint16 {
	opcode, exists := c.Ops[name]
	if !exists {
		opcode = 0x00
	}
	return opcode
}

// Op executes an opcode with the given operands.
func (c *CPU) Op(opcode uint16, operands []uint16) error {
	fmt.Println(opcode, operands, c.Regs)
	switch opcode {
	// nop
	case 0x00:
	// cop (reg to copy to, reg to copy from)
	case 0x01:
		c.Regs[operands[0]].Set(
			c.Regs[operands[1]].Get(),
		)
	// cpl (reg to copy to, value to copy)
	case 0x02:
		c.Regs[operands[0]].Set(operands[1])
	// str (reg with addr, reg with value)
	case 0x03:
		c.Mem.Set(
			c.Regs[operands[0]].Get(),
			c.Regs[operands[1]].Get(),
		)
	// ldr (reg to load to, reg with addr)
	case 0x04:
		c.Regs[operands[0]].Set(
			c.Mem.Get(
				c.Regs[operands[1]].Get(),
			),
		)
	}
	return nil
}
