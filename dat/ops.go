package dat

var (
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
		"psh": 0x14,
		"pop": 0x15,
		"ret": 0x16,
		"cal": 0x17,
		"cmp": 0x18,
		"cle": 0x19,
		"cln": 0x1a,
		"gto": 0x1b,
		"gte": 0x1c,
		"gtn": 0x1d,
	}

	// OpCodeToName is the reverse of OpNameToCode.
	// It is created from it at runtime.
	OpCodeToName = make(map[uint16]string)

	// OpNameToSize maps names to the instruction size.
	OpNameToSize = map[string]uint16{
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
		"psh": 1,
		"pop": 1,
		"ret": 0,
		"cal": 1,
		"cmp": 2,
		"cle": 1,
		"cln": 1,
		"gto": 1,
		"gte": 1,
		"gtn": 1,
	}
)
