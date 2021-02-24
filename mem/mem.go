// Package mem implements memory devices.
package mem

// AddressSpace maps 16-bit addresses to 16-bit values.
type AddressSpace map[uint16]uint16

// RAM represents 128K of word-based memory (64K addresses).
type RAM struct {
	Mem           AddressSpace
	VGAOffset     uint16
	VGAWidth      int
	VGAHeight     int
	StackMin      uint16
	StackMax      uint16
	ProgramOffset uint16
	HeapOffset    uint16
}

// NewRAM returns a pointer to a newly initialized RAM.
func NewRAM(a AddressSpace, vw int, vh int) *RAM {
	vgaSize := uint16(vh * vw)
	return &RAM{
		Mem:           a,
		VGAOffset:     0,
		VGAHeight:     vh,
		VGAWidth:      vw,
		StackMin:      vgaSize + 1,
		StackMax:      vgaSize + 303,
		ProgramOffset: vgaSize + 304,
		HeapOffset:    vgaSize + 304,
	}
}

// Get gets the value stored at a specified address.
func (m *RAM) Get(address uint16) uint16 {
	value, exists := m.Mem[address]
	if !exists {
		m.Mem[address] = 0
		value = 0
	}
	return value
}

// Set sets the specified address to the specified value.
func (m *RAM) Set(address uint16, value uint16) {
	m.Mem[address] = value
}
