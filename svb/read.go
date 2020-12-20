package svb

import (
	"fmt"
	"github.com/tteeoo/svc/dat"
)

// ParseBinary deserializes an SVB.
func ParseBinary(b []byte) (SVB, error) {
	if len(b) < 7 {
		return SVB{}, fmt.Errorf("file is invalid (smaller than 7 bytes)")
	}

	// Compensate for 4 byte magic number.
	r := b[4:]

	// Parse out constants
	consts := []Constant{}
	preSubs := []byte{}
	a, value := uint16(0), false
	for i := 0; i < len(r); i += 2 {
		if BytesToUint(r[i:i+2]) == 0xffff {
			preSubs = r[i:]
			break
		}
		if value {
			consts = append(consts, Constant{
				Address: a,
				Value:   BytesToUint(r[i : i+2]),
			})
			value = false
		} else {
			a = BytesToUint(r[i : i+2])
			value = true
		}
	}

	// Parse out subroutines
	doubled := []uint16{}
	for i := 0; i < len(preSubs); i += 2 {
		doubled = append(doubled, BytesToUint(preSubs[i:i+2]))
	}
	subs := []Subroutine{}
	idx := 0
	for idx < len(doubled) {

		// Initialize subroutine data
		a := doubled[idx+1]
		ints := []Instruction{}

		// Parse out instructions
		opIdx := idx + 2
		for {

			// Get opcode and opname
			op := doubled[opIdx]
			size := 0
			opName, exists := dat.OpCodeToName[op]
			if exists {
				size = dat.OpNameToSize[opName]
			} else {
				return SVB{}, fmt.Errorf("svb is invalid (opcode 0x%x does not exist)", op)
			}

			// Parse out operands
			operands := []uint16{}
			for _, and := range doubled[opIdx+1 : opIdx+size+1] {
				operands = append(operands, and)
			}

			ints = append(ints, Instruction{
				Name:     opName,
				Opcode:   op,
				Operands: operands,
				Size:     size + 1,
			})

			// Exit if opIdx is 0xffff, else continue parsing next int
			if (len(doubled) <= opIdx+size+1) || (doubled[opIdx+size+1] == 0xffff) {
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

		idx += size + 2
	}

	return SVB{
		Subroutines: subs,
		Constants:   consts,
		MainAddress: subs[len(subs)-1].Address,
	}, nil
}
