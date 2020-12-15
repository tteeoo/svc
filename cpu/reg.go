package cpu

import (
	"fmt"
)

// Register represents a CPU register for storing one word.
type Register struct {
	// State stores the value of the Register.
	State uint16
	// Num is the Number of the Register.
	Num uint16
}

// NewRegister returns a newly initialized Register, taking its Num.
func NewRegister(Num uint16) *Register {
	return &Register{Num: Num}
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
	return fmt.Sprintf("<%d:%x>", r.Num, r.State)
}
