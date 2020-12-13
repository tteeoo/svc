package mem

// GenericMemoryDevice is a basic implementation of a MemoryDevice.
type GenericMemoryDevice struct {
	mem addressSpace
}

// NewGenericMemoryDevice returns a pointer to a newly initialized GenericMemoryDevice.
func NewGenericMemoryDevice() *GenericMemoryDevice {
	return &GenericMemoryDevice{mem: addressSpace{}}
}

// Get gets the value stored at a specified address.
func (m *GenericMemoryDevice) Get(address uint16) uint16 {
	value, exists := m.mem[address]
	if !exists {
		m.mem[address] = 0x00
		value = 0x00
	}
	return value
}

// Set sets the specified address to the specified value.
func (m *GenericMemoryDevice) Set(address uint16, value uint16) {
	m.mem[address] = value
}
