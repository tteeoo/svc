// Package mem implements memory devices.
package mem

// AddressSpace maps 16-bit addresses to 16-bit values.
type AddressSpace map[uint16]uint16

// TextBuffer represents a specific part of memory used to store characters that
// would display on a screen.
type TextBuffer [25][80]uint16

// MemoryDevice is an interface that wraps methods to interact with a basic
// random-access memory device.
type MemoryDevice interface {
	GetTextBuffer() *TextBuffer
	GetAddressSpace(uint16) (*AddressSpace, error)
	NewAddressSpace() (uint16, *AddressSpace)
}
