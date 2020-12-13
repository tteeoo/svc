package cpu

import (
	"github.com/tteeoo/svc/mem"
)

// GenericCPU is a basic implementation of a CPU.
type GenericCPU struct {
	// memoryDevice is the memory device used by the CPU.
	memoryDevice mem.MemoryDevice
}

// NewGenericCPU returns a pointer to a newly initialized GenericCPU.
func NewGenericCPU() *GenericCPU {
	return &GenericCPU{
		memoryDevice: mem.NewGenericMemoryDevice(),
	}
}

// GetMemoryDevice returns the GenericCPU's memory device.
func (c *GenericCPU) GetMemoryDevices() mem.MemoryDevice {
	return c.memoryDevice
}

// Op executes an opcode with the given operands.
func (c *GenericCPU) Op(opcode uint16, operands []uint16) {
	switch opcode {
	case 0x00:
	}
}
