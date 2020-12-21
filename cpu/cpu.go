// Package cpu implements structs and functions for
// representing and operating the virtual CPU.
package cpu

import (
	"os"
	"fmt"
	"github.com/tteeoo/svc/dat"
	"github.com/tteeoo/svc/mem"
	"github.com/tteeoo/svc/vga"
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
func NewCPU(m *mem.RAM, v *vga.VGA) *CPU {

	// Create registers
	regs := make(map[uint16]*Register)
	for i := 0; i < dat.GPRNum; i++ {
		regs[uint16(i)] = NewRegister()
	}
	for _, i := range []string{"ex", "ac", "sp", "pc"} {
		regs[dat.RegNamesToNum[i]] = NewRegister()
	}
	regs[dat.RegNamesToNum["sp"]].Set(dat.StackOffset)

	return &CPU{
		Mem:  m,
		VGA:  v,
		Regs: regs,
	}
}

// Run starts execution at the given memory address.
func (c *CPU) Run(address uint16) {

	// Push exit address onto stack
	sp := dat.RegNamesToNum["sp"]
	c.Regs[sp].Set(c.Regs[sp].Get() - 1)
	c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]].Get(), 0xffff)

	c.Regs[dat.RegNamesToNum["pc"]].Set(address)
	for {
		pc := c.Regs[dat.RegNamesToNum["pc"]].Get()

		// Exit if pc is the last address
		if pc == 0xffff {
			os.Exit(0)
		}

		op := c.Mem.Get(pc)
		name := dat.OpCodeToName[op]
		size := dat.OpNameToSize[name]

		operands := make([]uint16, size)
		for i := 0; i < size; i++ {
			operands[i] = c.Mem.Get(pc+uint16(1+i))
		}

		c.Op(op, operands)
		c.Regs[dat.RegNamesToNum["pc"]].Set(pc+uint16(1+size))
	}
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
func (c *CPU) Op(opcode uint16, operands []uint16) {
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
	// psh (reg with value)
	case 0x14:
		c.Regs[dat.RegNamesToNum["sp"]].Set(
			c.Regs[dat.RegNamesToNum["sp"]].Get() - 1,
		)
		c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]].Get(), operands[0])
	// pop (reg to store in)
	case 0x15:
		c.Regs[operands[0]].Set(
			c.Regs[dat.RegNamesToNum["sp"]].Get(),
		)
		c.Regs[dat.RegNamesToNum["sp"]].Set(
			c.Regs[dat.RegNamesToNum["sp"]].Get() + 1,
		)
	// ret
	case 0x16:
		c.Regs[dat.RegNamesToNum["pc"]].Set(
			c.Regs[dat.RegNamesToNum["sp"]].Get(),
		)
		c.Regs[dat.RegNamesToNum["sp"]].Set(
			c.Regs[dat.RegNamesToNum["sp"]].Get() + 1,
		)
		// TODO:
		// cal
		// cmp
		// cle
		// cln
		// jmp
		// cle
		// cln
	}
	// fmt.Println(c)
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
