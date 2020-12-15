// Package mem implements memory devices.
package mem

import (
	"fmt"
)

// AddressSpace maps 16-bit addresses to 16-bit values.
type AddressSpace map[uint16]uint16

// RAM represents 128K of word-based memory (64K addresses).
type RAM struct {
	Mem AddressSpace
}

// NewRAM returns a pointer to a newly initialized RAM.
func NewRAM() *RAM {
	return &RAM{Mem: AddressSpace{}}
}

// Get gets the value stored at a specified address.
func (m *RAM) Get(address uint16) uint16 {
	value, exists := m.Mem[address]
	if !exists {
		// Because actual memory would have an arbitrary initial value, this value can be anything :)
		m.Mem[address] = 0x1337
		value = 0x1337
	}
	return value
}

// Set sets the specified address to the specified value.
func (m *RAM) Set(address uint16, value uint16) {
	m.Mem[address] = value
}

// String returns a string representation of RAM.
func (m *RAM) String() string {
	addrs := make([]uint16, len(m.Mem))
	i := 0
	for k := range m.Mem {
		addrs[i] = k
		i++
	}
	out := ""
	for _, j := range addrs {
		out += fmt.Sprintf("\n%x:%x", j, m.Mem[j])
	}
	if len(out) == 0 {
		return "no memory allocated"
	}
	return out[1:]
}
