package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/svb"
)

// TODO: color options
func color(s string, c string) string {
	return "\033[" + c + "m" + s + "\033[0m"
}

// TODO: util package
// parseHex will take a hex-representing string and parse it to a uint16.
func parseHex(s string) (uint16, error) {
	// Check the length
	if len(s) > 4 {
		return 0, fmt.Errorf("hex value \"%s\" is too large", s)
	}
	// Pad and decode
	b, err := hex.DecodeString(fmt.Sprintf("%0*s", 4, s))
	if err != nil {
		return 0, fmt.Errorf("cannot decode hex value \"%s\": %s", s, err)
	}
	return svb.BytesToUint(b), nil
}

func addressToSection(c *cpu.CPU, a uint16) string {
	if a < c.Mem.StackMin {
		return "text"
	} else if a < c.Mem.StackMax+1 {
		return "stak"
	} else if a < c.Mem.HeapOffset {
		return "prog"
	} else {
		return "heap"
	}
}

func sectionToColor(s string) string {
	var ansic string
	switch s {
	case "text":
		ansic = "35;1"
	case "stak":
		ansic = "36;1"
	case "prog":
		ansic = "32;1"
	case "heap":
		ansic = "33;1"
	}
	return ansic
}
