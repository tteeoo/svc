package svb

import (
	"bytes"
	"encoding/binary"
)

// Bytes serializes an SVB.
func (s SVB) Bytes() []byte {

	// Calculate length
	l := len(s.Constants) + len(s.Subroutines)
	for _, i := range s.Subroutines {
		l += i.Size
	}

	// Add constants
	u := make([]uint16, l)
	for i, c := range s.Constants {
		u[i] = c.Value
	}

	// Add subroutines
	i := len(s.Constants)
	for _, sub := range s.Subroutines {
		u[i] = 0xffff
		opi := i + 1
		for _, op := range sub.Instructions {
			u[opi] = op.Opcode
			opi++
			for _, and := range op.Operands {
				u[opi] = and
				opi++
			}
		}
		i += sub.Size + 1
	}

	// Convert uint16 to binary
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return append([]byte("SVCB"), buf.Bytes()...)
}
