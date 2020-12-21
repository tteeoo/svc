// Package cpu implements structs and functions for
// representing and operating the virtual CPU.
package cpu

import (
	"fmt"
	"github.com/tteeoo/svc/dat"
	"github.com/tteeoo/svc/mem"
	"github.com/tteeoo/svc/vga"
	"os"
)

// CPU is a basic implementation of a CPU.
type CPU struct {
	// Mem is the main memory device used by the CPU.
	Mem *mem.RAM
	// VGA is the main video device used by the CPU.
	VGA *vga.VGA
	// Regs maps numbers to regsiter values.
	Regs map[uint16]uint16
}

// NewCPU returns a pointer to a newly initialized CPU.
func NewCPU(m *mem.RAM, v *vga.VGA) *CPU {

	// Create registers
	regs := make(map[uint16]uint16)
	for i := 0; i < dat.GPRNum; i++ {
		regs[uint16(i)] = 0
	}
	for _, i := range []string{"ex", "ac", "sp", "pc"} {
		regs[dat.RegNamesToNum[i]] = 0
	}
	regs[dat.RegNamesToNum["sp"]] = dat.StackOffset

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
	c.Regs[sp]--
	c.Mem.Set(c.Regs[sp], 0xffff)

	c.Regs[dat.RegNamesToNum["pc"]] = address
	for {
		pc := c.Regs[dat.RegNamesToNum["pc"]]

		// Exit if pc is the last address
		if pc == 0xffff {
			os.Exit(0)
		}

		op := c.Mem.Get(pc)
		name := dat.OpCodeToName[op]
		size := dat.OpNameToSize[name]

		operands := make([]uint16, size)
		for i := 0; i < size; i++ {
			operands[i] = c.Mem.Get(pc + uint16(1+i))
		}

		c.Regs[dat.RegNamesToNum["pc"]] += uint16(1 + size)

		c.Op(op, operands)
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
		c.Regs[operands[0]] = c.Regs[operands[1]]
	// cpl (reg to copy to, value to copy)
	case 0x02:
		c.Regs[operands[0]] = operands[1]
	// str (reg with addr, reg with value)
	case 0x03:
		c.Mem.Set(
			c.Regs[operands[0]],
			c.Regs[operands[1]],
		)
	// ldr (reg to load to, reg with addr)
	case 0x04:
		c.Regs[operands[0]] = c.Mem.Get(c.Regs[operands[1]])
	// add (reg with value)
	case 0x05:
		c.Regs[dat.RegNamesToNum["ac"]] = c.Regs[dat.RegNamesToNum["ac"]] + c.Regs[operands[0]]
	// sub (reg with value)
	case 0x06:
		c.Regs[dat.RegNamesToNum["ac"]] = c.Regs[dat.RegNamesToNum["ac"]] + (^c.Regs[operands[0]] + 1)
	// twc (reg to twc)
	case 0x07:
		c.Regs[operands[0]] = ^c.Regs[operands[0]] + 1
	// inc (reg to inc)
	case 0x08:
		c.Regs[operands[0]]++
	// dec (reg to dec)
	case 0x09:
		c.Regs[operands[0]] += (^uint16(1) + 1)
	// mul (reg with value)
	case 0x0a:
		c.Regs[dat.RegNamesToNum["ac"]] *= c.Regs[operands[0]]
	// div (reg with value)
	case 0x0b:
		c.Regs[dat.RegNamesToNum["ex"]] = c.Regs[dat.RegNamesToNum["ac"]] % c.Regs[operands[0]]
		c.Regs[dat.RegNamesToNum["ac"]] /= c.Regs[operands[0]]
	// dvc (reg with value)
	case 0x0c:
		a := c.Regs[dat.RegNamesToNum["ac"]]
		b := c.Regs[operands[0]]
		x, y := a, b
		aSign, bSign := a>>15, b>>15
		same := aSign == bSign
		if aSign == 1 {
			x = ^x + 1
		}
		if bSign == 1 {
			y = ^y + 1
		}
		c.Regs[dat.RegNamesToNum["ex"]] = x % y
		if same {
			c.Regs[dat.RegNamesToNum["ac"]] = x / y
		} else {
			c.Regs[dat.RegNamesToNum["ac"]] = ^(x / y) + 1
		}
	// xor (reg with value)
	case 0x0d:
		c.Regs[dat.RegNamesToNum["ac"]] ^= c.Regs[operands[0]]
	// and (reg with value)
	case 0x0e:
		c.Regs[dat.RegNamesToNum["ac"]] &= c.Regs[operands[0]]
	// orr (reg with value)
	case 0x0f:
		c.Regs[dat.RegNamesToNum["ac"]] |= c.Regs[operands[0]]
	// not (reg to invert)
	case 0x10:
		c.Regs[operands[0]] = ^c.Regs[operands[0]]
	// shr (reg to shift, amount to shift)
	case 0x11:
		c.Regs[operands[0]] >>= operands[1]
	// shl (reg to shift, amount to shift)
	case 0x12:
		c.Regs[operands[0]] <<= operands[1]
	// vga
	case 0x13:
		c.VGA.TextDraw()
	// psh (reg with value)
	case 0x14:
		c.Regs[dat.RegNamesToNum["sp"]]--
		c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]], operands[0])
	// pop (reg to store in)
	case 0x15:
		c.Regs[operands[0]] = c.Mem.Get(c.Regs[dat.RegNamesToNum["sp"]])
		c.Regs[dat.RegNamesToNum["sp"]]++
	// ret
	case 0x16:
		c.Regs[dat.RegNamesToNum["pc"]] = c.Mem.Get(c.Regs[dat.RegNamesToNum["sp"]])
		c.Regs[dat.RegNamesToNum["sp"]]++
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

	out := "Memory:\n"
	out += c.Mem.String()

	out += "\nRegisters:"
	for i := 0; i < dat.GPRNum; i++ {
		out += fmt.Sprintf("\n%d:%x", i, c.Regs[uint16(i)])
	}

	out += fmt.Sprintf("\nex:%x", c.Regs[dat.RegNamesToNum["ex"]])
	out += fmt.Sprintf("\nac:%x", c.Regs[dat.RegNamesToNum["ac"]])
	out += fmt.Sprintf("\nsp:%x", c.Regs[dat.RegNamesToNum["sp"]])
	out += fmt.Sprintf("\npc:%x", c.Regs[dat.RegNamesToNum["pc"]])

	return out
}
