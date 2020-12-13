package mem

import (
	"fmt"
)

// GenericMemoryDevice is a basic implementation of a MemoryDevice.
// It has one 80x25 text-mode display buffer.
// It can have multiple "memory banks" which allow it to store more than just 64K.
type GenericMemoryDevice struct {
	tb    TextBuffer
	banks map[uint16]*AddressSpace
}

// NewGenericMemoryDevice returns a pointer to a newly initialized GenericMemoryDevice.
func NewGenericMemoryDevice() *GenericMemoryDevice {
	return &GenericMemoryDevice{tb: TextBuffer{}, banks: make(map[uint16]*AddressSpace)}
}

// GetTextBuffer returns a pointer to the device's text buffer.
func (m *GenericMemoryDevice) GetTextBuffer() *TextBuffer {
	return &m.tb
}

// GetAddressSpace returns a pointer the address space of the specified bank.
// It returns an error if the address space does not exist.
func (m *GenericMemoryDevice) GetAddressSpace(key uint16) (*AddressSpace, error) {
	value, exists := m.banks[key]
	if !exists {
		return nil, fmt.Errorf("address space %x does not exist", key)
	}
	return value, nil
}

// NewAddressSpace creates a new memory bank and returns its address and
// a pointer to its address space.
func (m *GenericMemoryDevice) NewAddressSpace() (uint16, *AddressSpace) {
	key := uint16(len(m.banks))
	as := &AddressSpace{}
	m.banks[key] = as
	return key, as
}
