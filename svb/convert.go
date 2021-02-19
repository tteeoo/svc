package svb

import (
	"bytes"
	"encoding/binary"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/mem"
	"github.com/tteeoo/svc/util"
)

// LoadProgram takes the bytes of an SVB file and parses out
//   the new address space, main subroutine address, and program size.
func LoadProgram(c *cpu.CPU, b []byte) (mem.AddressSpace, uint16, uint16) {

	// []byte -> []uint16
	u := make([]uint16, len(b)/2)
	for i := 0; i < cap(u); i++ {
		u[i] = util.BytesToUint([]byte{b[i*2], b[(i*2)+1]})
	}

	// Extract headers (extensible)
	mainAddress := uint16(0)
	headerIndex := 0
	for i := 0; i < len(u); i++ {
		if u[i] == 0xffff {
			headerIndex = i
			break
		}
		switch i {
		case 0:
			mainAddress = u[i]
		}
	}

	// []uint16 -> address space
	as := make(mem.AddressSpace)
	for i, j := range u[headerIndex+1:] {
		as[c.Mem.ProgramOffset+uint16(i)] = j
	}

	return as, mainAddress, uint16(len(u[headerIndex+1:]))
}

// Bytes serializes an SVB.
func (s SVB) Bytes() []byte {

	// Add constants
	u := make([]uint16, s.Size()+2)
	u[0] = s.MainAddress
	u[1] = 0xffff
	for i, c := range s.Constants {
		u[i+2] = c.Value
	}

	// Add subroutines
	i := uint16(len(s.Constants)) + 2
	for _, sub := range s.Subroutines {
		for _, op := range sub.Instructions {
			u[i] = op.Opcode
			i++
			for _, and := range op.Operands {
				u[i] = and
				i++
			}
		}
	}

	// Convert uint16 to binary
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return buf.Bytes()
}
