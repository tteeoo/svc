// Package mem implements memory devices.
package mem

// addressSpace maps 16-bit addresses to 16-bit values.
type addressSpace map[uint16]uint16

// MemoryDevice is an interface that wraps methods to interact with a basic
// random-access memory device.
type MemoryDevice interface {
	Get(uint16) uint16
	Set(uint16, uint16)
}
