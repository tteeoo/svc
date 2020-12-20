package svb

import (
	"github.com/tteeoo/svc/mem"
)

// GetProgramMem returns an address space with the program loaded in it.
func (s SVB) GetProgramMem() mem.AddressSpace {

	a := make(mem.AddressSpace)

	for _, c := range s.Constants {
		a[c.Address] = c.Value
	}

	for _, s := range s.Subroutines {
		x := s.Address
		for _, i := range s.Instructions {
			a[x] = i.Opcode
			x++
			for _, and := range i.Operands {
				a[x] = and
				x++
			}
		}
	}

	return a
}
