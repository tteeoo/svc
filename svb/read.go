package svb

import (
	"fmt"
	"github.com/tteeoo/svc/dat"
)

// MainAddress is the address of the subroutine that should be main.
var MainAddress = uint16(0x900)

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
		name := ""
		ints := []Instruction{}
		if a == MainAddress {
			name = "main"
		}

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
				return SVB{}, fmt.Errorf("svb is invalid (opcode %x does not exist)", op)
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
				Size:     size,
			})

			// Exit if opIdx is 0xffff, else continue parsing next int
			if doubled[opIdx+size+1] == 0xffff {
				break
			} else {
				opIdx = opIdx + size + 1
			}
		}

		// Get size
		size := 0
		for _, i := range ints {
			size += i.Size
		}

		subs = append(subs, Subroutine{
			Name:         name,
			Address:      a,
			Instructions: ints,
			Size:         size,
		})

		idx = idx + size + 1
	}

	return SVB{
		Subroutines: subs,
		Constants:   consts,
	}, nil
}
