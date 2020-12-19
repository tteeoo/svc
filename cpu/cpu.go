// Package cpu implements structs and functions for
// representing and operating the virtual CPU.
package cpu

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
	"github.com/tteeoo/svc/vga"
	"github.com/tteeoo/svc/dat"
)

// CPU is a basic implementation of a CPU.
type CPU struct {
	// Mem is the main memory device used by the CPU.
	Mem *mem.RAM
	// VGA is the main video device used by the CPU.
	VGA *vga.VGA
	// Regs maps numbers to regsiters.
	Regs map[uint16]*Register
}

// NewCPU returns a pointer to a newly initialized CPU.
func NewCPU() *CPU {

	// Create registers
	regs := make(map[uint16]*Register)
	for i := 0; i < dat.GPRNum; i++ {
		regs[uint16(i)] = NewRegister()
	}
	for _, i := range []string{"ex", "ac", "sp", "pc"} {
		regs[dat.RegNamesToNum[i]] = NewRegister()
	}

	m := mem.NewRAM()
	return &CPU{
		Mem:  m,
		VGA:  vga.NewVGA(m, 0x00, 25, 80),
		Regs: regs,
	}
}

// GetMem returns the CPU's memory device.
func (c *CPU) GetMem() *mem.RAM {
	return c.Mem
}

// GetOp returns to opcode whose name is provided.
// Returns 0x00 (nop) if the name is not defined.
func (c *CPU) GetOp(name string) uint16 {
	opcode, exists := dat.OpNameToCode[name]
	if !exists {
		opcode = 0x00
	}
	return opcode
}

// Op executes an opcode with the given operands.
func (c *CPU) Op(opcode uint16, operands []uint16) error {
	// fmt.Printf("--> %x, %x\n", opcode, operands)
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
	// add (reg with value)
	case 0x05:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() + c.Regs[operands[0]].Get(),
		)
	// sub (reg with value)
	case 0x06:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() + (^c.Regs[operands[0]].Get() + 1),
		)
	// twc (reg to twc)
	case 0x07:
		c.Regs[operands[0]].Set(
			^c.Regs[operands[0]].Get() + 1,
		)
	// inc (reg to inc)
	case 0x08:
		c.Regs[operands[0]].Set(
			c.Regs[operands[0]].Get() + 1,
		)
	// dec (reg to dec)
	case 0x09:
		c.Regs[operands[0]].Set(
			c.Regs[operands[0]].Get() + (^uint16(1) + 1),
		)
	// mul (reg with value)
	case 0x0a:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() * c.Regs[operands[0]].Get(),
		)
	// div (reg with value)
	case 0x0b:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() % c.Regs[operands[0]].Get(),
		)
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() / c.Regs[operands[0]].Get(),
		)
	// dvc (reg with value)
	case 0x0c:
		a := c.Regs[dat.RegNamesToNum["ac"]].Get()
		b := c.Regs[operands[0]].Get()
		x, y := a, b
		aSign, bSign := a>>15, b>>15
		same := aSign == bSign
		if aSign == 1 {
			x = ^x + 1
		}
		if bSign == 1 {
			y = ^y + 1
		}
		c.Regs[dat.RegNamesToNum["ex"]].Set(x % y)
		if same {
			c.Regs[dat.RegNamesToNum["ac"]].Set(x / y)
		} else {
			c.Regs[dat.RegNamesToNum["ac"]].Set(^(x / y) + 1)
		}
	// xor (reg with value)
	case 0x0d:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() ^ c.Regs[operands[0]].Get(),
		)
	// and (reg with value)
	case 0x0e:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() & c.Regs[operands[0]].Get(),
		)
	// orr (reg with value)
	case 0x0f:
		c.Regs[dat.RegNamesToNum["ac"]].Set(
			c.Regs[dat.RegNamesToNum["ac"]].Get() | c.Regs[operands[0]].Get(),
		)
	// not (reg to invert)
	case 0x10:
		c.Regs[operands[0]].Set(
			^c.Regs[operands[0]].Get(),
		)
	// shr (reg to shift, amount to shift)
	case 0x11:
		c.Regs[operands[0]].Set(
			c.Regs[operands[0]].Get() >> operands[1],
		)
	// shl (reg to shift, amount to shift)
	case 0x12:
		c.Regs[operands[0]].Set(
			c.Regs[operands[0]].Get() << operands[1],
		)
	// vga
	case 0x13:
		c.VGA.TextDraw()
		// TODO:
		// psh
		// pop
		// jmp
		// cmp
		// jme
		// jne
	}
	// fmt.Println(c)
	return nil
}

// String returns a string representation of a CPU.
func (c *CPU) String() string {
	out := "Registers:"
	for i := 0; i < dat.GPRNum; i++ {
		out += fmt.Sprintf("\n%d:%s", i, c.Regs[uint16(i)])
	}
	out += fmt.Sprintf("\nex:%s", c.Regs[dat.RegNamesToNum["ex"]])
	out += fmt.Sprintf("\nac:%s", c.Regs[dat.RegNamesToNum["ac"]])
	out += fmt.Sprintf("\nsp:%s", c.Regs[dat.RegNamesToNum["sp"]])
	out += fmt.Sprintf("\npc:%s", c.Regs[dat.RegNamesToNum["pc"]])
	out += "\nMemory:\n"
	out += c.Mem.String()
	return out
}
