package dat

const (
	// GPRNum is the number of general purpose registers.
	GPRNum = 8
)

var (
	// RegNamesToNum maps register names to numbers.
	RegNamesToNum = map[string]uint16{
		"ra": 0,
		"rb": 1,
		"rc": 2,
		"rd": 3,
		"re": 4,
		"rf": 5,
		"rh": 6,
		"ri": 7,
		"ex": 8,
		"ac": 9,
		"sp": 10,
		"pc": 11,
		"bi": 12,
	}
)
