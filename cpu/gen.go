package cpu

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
)

// GenericCPU is a basic implementation of a CPU.
type GenericCPU struct {
	// mem is the memory device used by the CPU.
	mem mem.MemoryDevice
	// regs maps numbers to regsiters.
	regs map[uint16]Register
	// ops maps opcode names to opcodes.
	ops map[string]uint16
}

// NewGenericCPU returns a pointer to a newly initialized GenericCPU.
func NewGenericCPU() *GenericCPU {
	regs := make(map[uint16]Register)
	// general purpose registers
	for _, i := range []uint16{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07} {
		regs[i] = NewGPR(i)
	}
	// instruction pointer
	regs[0x08] = NewGPR(0x08)
	// accumulator
	regs[0x09] = NewNoWriteR(0x09)
	return &GenericCPU{
		mem:  mem.NewGenericMemoryDevice(),
		regs: regs,
		ops: map[string]uint16{
			"nop": 0x00,
			"cop": 0x01,
			"cpl": 0x02,
			"str": 0x03,
			"ldr": 0x04,
		},
	}
}

// GetMemoryDevice returns the GenericCPU's memory device.
func (c *GenericCPU) GetMemoryDevice() mem.MemoryDevice {
	return c.mem
}

// GetOp returns to opcode whose name is provided.
// Returns 0x00 (nop) if the name is not defined.
func (c *GenericCPU) GetOp(name string) uint16 {
	opcode, exists := c.ops[name]
	if !exists {
		opcode = 0x00
	}
	return opcode
}

// Op executes an opcode with the given operands.
func (c *GenericCPU) Op(opcode uint16, operands []uint16) error {
	fmt.Println(opcode, operands, c.regs)
	switch opcode {
	// nop
	case 0x00:
	// cop (reg to copy to, reg to copy from)
	case 0x01:
		c.regs[operands[0]].Set(
			c.regs[operands[1]].Get(),
		)
	// cpl (reg to copy to, value to copy)
	case 0x02:
		c.regs[operands[0]].Set(operands[1])
	// str (reg with addr, reg with value)
	case 0x03:
		c.mem.Set(
			c.regs[operands[0]].Get(),
			c.regs[operands[1]].Get(),
		)
	// ldr (reg to load to, reg with addr)
	case 0x04:
		c.regs[operands[0]].Set(
			c.mem.Get(
				c.regs[operands[1]].Get(),
			),
		)
	}
	return nil
}
