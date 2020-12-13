package cpu

import (
	"fmt"
)

// Register wraps methods for a CPU register.
type Register interface {
	Get() uint16
	Set(uint16)
}

// GPR implements a CPU's general purpose register.
type GPR struct {
	// state stores the value of the GPR.
	state uint16
	// num is the number of the GPR.
	num uint16
}

// NewGPR returns a pointer to a newly initialized GPR, taking its num.
func NewGPR(num uint16) *GPR {
	return &GPR{num: num}
}

// Get returns the value of the GPR.
func (r *GPR) Get() uint16 {
	return r.state
}

// Set sets the value of the GPR.
func (r *GPR) Set(value uint16) {
	r.state = value
}

// String returns a string representation of the GPR.
func (r *GPR) String() string {
	return fmt.Sprintf("<GPR%d:%d>", r.num, r.state)
}

// NoWriteR is a register that cannot be written to with .Set().
type NoWriteR struct {
	// state stores the value of the NoWriteR.
	state uint16
	// num is the number of the NoWriteR.
	num uint16
}

// NewNoWriteR returns a pointer to a newly initialized NoWriteR, taking its num.
func NewNoWriteR(num uint16) *NoWriteR {
	return &NoWriteR{num: num}
}

// Get returns the value of the GPR.
func (r *NoWriteR) Get() uint16 {
	return r.state
}

// Set sets the value of the GPR.
func (r *NoWriteR) Set(value uint16) {}

// String returns a string representation of the NoWriteR.
func (r *NoWriteR) String() string {
	return fmt.Sprintf("<NWR%d:%d>", r.num, r.state)
}
