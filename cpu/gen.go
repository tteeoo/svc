package cpu

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
)

// GenericCPU is a basic implementation of a CPU.
type GenericCPU struct {
	// ram is the memory device used by the CPU.
	ram mem.MemoryDevice
	// gprs contains the eight general-purpose registers.
	gprs [8]uint16
	// ops maps opcode names to opcodes.
	ops map[string]uint16
}

// NewGenericCPU returns a pointer to a newly initialized GenericCPU.
func NewGenericCPU() *GenericCPU {
	return &GenericCPU{
		ram: mem.NewGenericMemoryDevice(),
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
	return c.ram
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
	fmt.Println(opcode, operands, c.gprs)
	switch opcode {
	// nop
	case 0x00:

	// cop (reg to copy from, reg to copy to)
	case 0x01:
		if len(operands) != 2 {
			return fmt.Errorf("cop requires two operands, %d were given", len(operands))
		}
		c.gprs[operands[1]] = c.gprs[operands[0]]

	// cpl (value to copy, reg to copy to)
	case 0x02:
		if len(operands) != 2 {
			return fmt.Errorf("cpl requires two operands, %d were given", len(operands))
		}
		c.gprs[operands[1]] = operands[0]

	// str (reg with value, reg with bank, reg with addr)
	case 0x03:
		if len(operands) != 3 {
			return fmt.Errorf("str requires three operands, %d were given", len(operands))
		}
		// Get the memory bank
		bank := *c.ram.GetAddressSpace(c.gprs[operands[1]])
		// Set the address of the memory bank
		bank[c.gprs[operands[2]]] = c.gprs[operands[0]]

	// ldr (reg with bank, reg with addr, reg to load to)
	case 0x04:
		if len(operands) != 3 {
			return fmt.Errorf("ldr requires three operands, %d were given", len(operands))
		}
		// Get the memory bank
		bank := *c.ram.GetAddressSpace(operands[0])
		// Set the register
		c.gprs[operands[2]] = bank[operands[1]]

	}
	return nil
}
