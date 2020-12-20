package svb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Bytes serializes an SVB.
func (s SVB) Bytes() ([]byte, error) {

	// Check validity
	for _, i := range s.Constants {
		if i.Address == 0xffff {
			return nil, fmt.Errorf("cannot initialize memory address 0xffff in svb file")
		}
	}

	// Calculate length
	l := len(s.Constants)*2 + len(s.Subroutines)*2
	for _, i := range s.Subroutines {
		l += i.Size
	}

	// Add constants
	u := make([]uint16, l)
	i := 0
	for _, c := range s.Constants {
		u[i] = c.Address
		u[i+1] = c.Value
		i += 2
	}

	// Add subroutines
	for _, sub := range s.Subroutines {
		u[i] = 0xffff
		u[i+1] = sub.Address
		opi := i + 2
		for _, op := range sub.Instructions {
			u[opi] = op.Opcode
			for _, and := range op.Operands {
				opi++
				u[opi] = and
			}
		}
		i += sub.Size + 2
	}

	// Convert uint16 to binary
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return append([]byte("SVCB"), buf.Bytes()...), nil
}
