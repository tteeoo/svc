package svb

import (
	"fmt"
	"github.com/tteeoo/svc/dat"
)

// ParseBinary deserializes an SVB.
func ParseBinary(b []byte) (SVB, error) {
	if len(b) < 8 {
		return SVB{}, fmt.Errorf("file is invalid (smaller than 8 bytes)")
	}

	// Compensate for 4 byte magic number.
	if string(b[:4]) != "SVCB" {
		return SVB{}, fmt.Errorf("file is not svb format")
	}
	r := b[4:]

	// Convert bytes to uints
	doubled := make([]uint16, len(r)/2)
	for i := 0; i < cap(doubled); i++ {
		doubled[i] = BytesToUint(r[2*i : 2*i+2])
	}

	// Parse out constants
	consts := []Constant{}
	bSubs := []uint16{}
	a := dat.ProgramOffset
	for i := 0; i < len(r); i++ {
		if doubled[i] == 0xffff {
			bSubs = doubled[i:]
			break
		}
		consts = append(consts, Constant{
			Address: a,
			Value:   doubled[i],
		})
		a++
	}

	// Parse out subroutines
	subs := []Subroutine{}
	idx := uint16(0)
	for int(idx) < len(bSubs) {

		// Initialize subroutine data
		ints := []Instruction{}

		// Parse out instructions
		opIdx := idx + 1
		for {

			// Get opcode and opname
			op := bSubs[opIdx]
			opSize := uint16(0)
			opName, exists := dat.OpCodeToName[op]
			if exists {
				opSize = dat.OpNameToSize[opName]
			} else {
				return SVB{}, fmt.Errorf("svb is invalid (opcode 0x%x does not exist)", op)
			}

			// Parse out operands
			operands := []uint16{}
			for _, and := range bSubs[opIdx+1 : opIdx+opSize+1] {
				operands = append(operands, and)
			}

			ints = append(ints, Instruction{
				Name:     opName,
				Opcode:   op,
				Operands: operands,
			})

			// Exit if opIdx is 0xffff, else continue parsing next int
			if (len(bSubs) <= int(opIdx+opSize+1)) || (bSubs[opIdx+opSize+1] == 0xffff) {
				break
			} else {
				opIdx += opSize + 1
			}
		}

		currentSub := Subroutine{
			Address:      a,
			Instructions: ints,
		}
		subs = append(subs, currentSub)
		idx += currentSub.Size() + 1
		a += currentSub.Size()
	}

	return SVB{
		Subroutines: subs,
		Constants:   consts,
		MainAddress: subs[len(subs)-1].Address,
	}, nil
}
