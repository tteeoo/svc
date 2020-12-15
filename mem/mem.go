// Package mem implements memory devices.
package mem

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
		m.Mem[address] = 0x00
		value = 0x00
	}
	return value
}

// Set sets the specified address to the specified value.
func (m *RAM) Set(address uint16, value uint16) {
	m.Mem[address] = value
}
