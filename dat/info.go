// Package dat contains information about CPU opcodes and registers.
package dat

const (
	// GPRNum is the number of general purpose registers.
	GPRNum = 4
)

var (
	// RegNamesToNum maps register names to numbers.
	RegNamesToNum = map[string]uint16{
		"ex": 4,
		"ac": 5,
		"sp": 6,
		"pc": 7,
	}

	// OpNameToCode maps names to actual opcodes.
	OpNameToCode = map[string]uint16{
		"nop": 0x00,
		"cop": 0x01,
		"cpl": 0x02,
		"str": 0x03,
		"ldr": 0x04,
		"add": 0x05,
		"sub": 0x06,
		"twc": 0x07,
		"inc": 0x08,
		"dec": 0x09,
		"mul": 0x0a,
		"div": 0x0b,
		"dvc": 0x0c,
		"xor": 0x0d,
		"and": 0x0e,
		"orr": 0x0f,
		"not": 0x10,
		"shr": 0x11,
		"shl": 0x12,
		"vga": 0x13,
	}

	// OpCodeToName is the reverse of OpNameToCode.
	// It is created from it at runtime.
	OpCodeToName = make(map[uint16]string)

	// OpNameToSize maps names to the instruction size.
	OpNameToSize = map[string]int{
		"nop": 0,
		"cop": 2,
		"cpl": 2,
		"str": 2,
		"ldr": 2,
		"add": 1,
		"sub": 1,
		"twc": 1,
		"inc": 1,
		"dec": 1,
		"mul": 1,
		"div": 1,
		"dvc": 1,
		"xor": 1,
		"and": 1,
		"orr": 1,
		"not": 1,
		"shr": 2,
		"shl": 2,
		"vga": 0,
	}
)

func init() {
	// Create the CodeToName map.
	for k, v := range OpNameToCode {
		OpCodeToName[v] = k
	}
}
