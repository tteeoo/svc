// Package cpu implements structs and functions for
// representing and operating the virtual CPU.
package cpu

import (
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
	for _, n := range dat.RegNamesToNum {
		regs[n] = 0
	}
	regs[dat.RegNamesToNum["sp"]] = m.StackMax

	return &CPU{
		Mem:  m,
		VGA:  v,
		Regs: regs,
	}
}

// Run starts execution at the given memory address.
func (c *CPU) Run(address uint16) {

	// Put command-line args into heap
	var l uint16
	if len(os.Args) > 2 {
		i := c.Mem.HeapOffset
		for _, str := range os.Args[2:] {
			for _, char := range str {
				c.Mem.Set(i, uint16(char))
				i++
			}
			c.Mem.Set(i, 0)
			i++
		}
		l = uint16(i - c.Mem.HeapOffset)
	}

	// Load heap information
	c.Mem.Set(0xffff, l)
	c.Mem.Set(0xfffe, c.Mem.HeapOffset)

	// Push exit address onto stack
	sp := dat.RegNamesToNum["sp"]
	c.Mem.Set(c.Regs[sp], 0xffff)

	// Set the program counter
	c.Regs[dat.RegNamesToNum["pc"]] = address

	// Enter the execution loop
	for {
		pc := c.Regs[dat.RegNamesToNum["pc"]]

		// Exit if pc is the last address
		if pc == 0xffff {
			os.Exit(0)
		}

		// Get instruction
		op := c.Mem.Get(pc)
		name := dat.OpCodeToName[op]
		size := dat.OpNameToSize[name]
		operands := make([]uint16, size)
		for i := uint16(0); i < size; i++ {
			operands[i] = c.Mem.Get(pc + (1 + i))
		}

		// Set the lc register
		c.Regs[dat.RegNamesToNum["lc"]] = pc

		// Increase program counter
		c.Regs[dat.RegNamesToNum["pc"]] += uint16(1 + size)

		// Execute instruction
		c.Op(op, operands)
	}
}

// Op executes an opcode with the given operands.
func (c *CPU) Op(opcode uint16, operands []uint16) {
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
		c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]], c.Regs[operands[0]])
	// pop (reg to store in)
	case 0x15:
		c.Regs[operands[0]] = c.Mem.Get(c.Regs[dat.RegNamesToNum["sp"]])
		c.Regs[dat.RegNamesToNum["sp"]]++
	// ret
	case 0x16:
		c.Regs[dat.RegNamesToNum["pc"]] = c.Mem.Get(c.Regs[dat.RegNamesToNum["sp"]])
		c.Regs[dat.RegNamesToNum["sp"]]++
	// cal (address to jump to)
	case 0x17:
		c.Regs[dat.RegNamesToNum["sp"]]--
		c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]], c.Regs[dat.RegNamesToNum["pc"]])
		c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
	// cmp (register, register)
	case 0x18:
		if c.Regs[operands[0]] == c.Regs[operands[1]] {
			c.Regs[dat.RegNamesToNum["bi"]] = 0xffff
		} else {
			c.Regs[dat.RegNamesToNum["bi"]] = 0xfffe
		}
	// cle (address to jump to)
	case 0x19:
		if c.Regs[dat.RegNamesToNum["bi"]] == 0xffff {
			c.Regs[dat.RegNamesToNum["sp"]]--
			c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]], c.Regs[dat.RegNamesToNum["pc"]])
			c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
		}
	// cln (address to jump to)
	case 0x1a:
		if c.Regs[dat.RegNamesToNum["bi"]] == 0xfffe {
			c.Regs[dat.RegNamesToNum["sp"]]--
			c.Mem.Set(c.Regs[dat.RegNamesToNum["sp"]], c.Regs[dat.RegNamesToNum["pc"]])
			c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
		}
	// gto (register holding address to jump to)
	case 0x1b:
		c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
	// gte (register holding address to jump to)
	case 0x1c:
		if c.Regs[dat.RegNamesToNum["bi"]] == 0xffff {
			c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
		}
	// gtn (register holding address to jump to)
	case 0x1d:
		if c.Regs[dat.RegNamesToNum["bi"]] == 0xfffe {
			c.Regs[dat.RegNamesToNum["pc"]] = operands[0]
		}
	// sth (reg with addr, reg with value)
	case 0x1e:
		c.Mem.Set(
			c.Regs[operands[0]]+c.Mem.HeapOffset,
			c.Regs[operands[1]],
		)
	// ldh (reg to load to, reg with addr)
	case 0x1f:
		c.Regs[operands[0]] = c.Mem.Get(c.Regs[operands[1]] + c.Mem.HeapOffset)
	}
}
