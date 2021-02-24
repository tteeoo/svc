// Package util contains utility functions.
package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tteeoo/svc/cpu"
)

// UintToBytes converts a uint16 to two bytes.
func UintToBytes(u uint16) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return buf.Bytes()
}

// BytesToUint converts two bytes to a uint16.
func BytesToUint(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// ParseHex takes a hex-representing string and parse it to a uint16.
func ParseHex(s string) (uint16, error) {
	// Check the length
	if len(s) > 4 {
		return 0, fmt.Errorf("hex value \"%s\" is too large", s)
	}
	// Pad and decode
	b, err := hex.DecodeString(fmt.Sprintf("%0*s", 4, s))
	if err != nil {
		return 0, fmt.Errorf("cannot decode hex value \"%s\": %s", s, err)
	}
	return BytesToUint(b), nil
}

// Color creates an ANSI code colored string.
func Color(s string, c string) string {
	return "\033[" + c + "m" + s + "\033[0m"
}

// AddressToSection takes memory address, outputs section string.
func AddressToSection(c *cpu.CPU, a uint16) string {
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

// SectionToColor takes section string, outputs color code.
func SectionToColor(s string) string {
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
