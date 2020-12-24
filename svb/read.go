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
	fmt.Printf("%x\n", doubled)

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
	idx := 0
	for idx < len(bSubs) {

		// Initialize subroutine data
		ints := []Instruction{}

		// Parse out instructions
		opIdx := idx + 1
		for {

			// Get opcode and opname
			op := bSubs[opIdx]
			size := 0
			opName, exists := dat.OpCodeToName[op]
			if exists {
				size = dat.OpNameToSize[opName]
			} else {
				return SVB{}, fmt.Errorf("svb is invalid (opcode 0x%x does not exist)", op)
			}

			// Parse out operands
			operands := []uint16{}
			for _, and := range bSubs[opIdx+1 : opIdx+size+1] {
				operands = append(operands, and)
			}

			ints = append(ints, Instruction{
				Name:     opName,
				Opcode:   op,
				Operands: operands,
				Size:     size + 1,
			})

			// Exit if opIdx is 0xffff, else continue parsing next int
			if (len(bSubs) <= opIdx+size+1) || (bSubs[opIdx+size+1] == 0xffff) {
				break
			} else {
				opIdx += size + 1
			}
		}

		// Get size
		size := 0
		for _, i := range ints {
			size += i.Size
		}

		subs = append(subs, Subroutine{
			Address:      a,
			Instructions: ints,
			Size:         size,
		})
		idx += size + 1
		a += uint16(idx)
	}

	return SVB{
		Subroutines: subs,
		Constants:   consts,
		MainAddress: subs[len(subs)-1].Address,
	}, nil
}
