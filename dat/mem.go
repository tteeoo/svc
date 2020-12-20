package dat

var (
	// VGAHeight is the height of the VGA text buffer.
	VGAHeight = 25

	// VGAWidth is the width of the VGA text buffer.
	VGAWidth = 80

	// VGAOffset is the offset that the VGA text buffer is stored at in memory.
	VGAOffset = uint16(0x0)

	// ProgramOffset is the offset that the program is stored at in memory.
	ProgramOffset = uint16(0x900)

	// StackOffset is the offset that the stack is stored at in memory.
	// Note that the stack counts down.
	StackOffset = uint16(0x8ff)
)
