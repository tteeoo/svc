package dat

const (
	// GPRNum is the number of general purpose registers.
	GPRNum = 4
)

var (
	// RegNamesToNum maps register names to numbers.
	RegNamesToNum = map[string]uint16{
		"aa": 0,
		"bb": 1,
		"cc": 2,
		"dd": 3,
		"ex": 4,
		"ac": 5,
		"sp": 6,
		"pc": 7,
		"bi": 8,
	}
)
