package cpu

import (
	"fmt"
)

// Register represents a CPU register for storing one word.
type Register struct {
	// State stores the value of the Register.
	State uint16
}

// NewRegister returns a pointer to a newly initialized Register.
func NewRegister() *Register {
	return &Register{}
}

// Get returns the value of the Register.
func (r *Register) Get() uint16 {
	return r.State
}

// Set sets the value of the Register.
func (r *Register) Set(value uint16) {
	r.State = value
}

// String returns a string representation of the Register.
func (r *Register) String() string {
	return fmt.Sprintf("%x", r.State)
}
