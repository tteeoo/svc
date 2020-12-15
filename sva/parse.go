package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"strconv"
	"strings"
)

// parse parses an input file into memory and opcodes.
func parse(b []byte) (map[uint16]uint16, [][]uint16, error) {

	var ops [][]uint16
	var address uint16 = 0x00
	mem := make(map[uint16]uint16)
	vars := make(map[string]uint16)

	// CPU for translating opcodes
	c := cpu.NewCPU()

	// Iterate lines
	split := strings.Split(string(b), "\n")
	for _, line := range split {

		// Parse out comments
		noComments := ""
		for _, char := range line {
			if char == ';' {
				break
			}
			noComments += string(char)
		}

		// Tokenize
		badSplitLine := strings.Split(strings.Replace(noComments, "\t", "", -1), " ")

		// Remove empty strings
		var splitLine []string
		for _, str := range badSplitLine {
			if str != "" {
				splitLine = append(splitLine, str)
			}
		}

		// Handle mem
		if (len(splitLine) == 3) && (splitLine[1] == "=") {
			if len(splitLine[2]) > 2 {

				// Handle a string
				if (splitLine[2][0] == byte('"')) && (splitLine[2][len(splitLine[2])-1] == byte('"')) {
					vars[splitLine[0]] = address
					for _, char := range splitLine[2][1 : len(splitLine[2])-1] {
						mem[address] = uint16(char)
						address++
					}
					continue

				} else if splitLine[2][1] == 'x' {
					// Handle a hex value
					val := splitLine[2][2:]
					if len(val) == 2 {
						val = "00" + val
					}
					vars[splitLine[0]] = address
					b, err := hex.DecodeString(val)
					if err != nil {
						return nil, nil, err
					}
					mem[address] = binary.BigEndian.Uint16(b)
					address++
					continue
				}
			}

			// Handle an int
			i, err := strconv.Atoi(splitLine[2])
			if err != nil {
				return nil, nil, err
			}
			vars[splitLine[0]] = address
			mem[address] = uint16(i)
			address++

		} else if len(splitLine) > 0 {
			// Handle instruction
			op := make([]uint16, len(splitLine))
			op[0] = c.GetOp(splitLine[0])

			for i, j := range splitLine[1:] {
				if (len(j) > 2) && (j[1] == 'x') {
					// Handle hex number
					val := j[2:]
					if len(val) == 2 {
						val = "00" + val
					}
					b, err := hex.DecodeString(val)
					if err != nil {
						return nil, nil, err
					}
					op[i+1] = binary.BigEndian.Uint16(b)

				} else if (len(j) > 2) && (j[0] == '[') && (j[len(j)-1] == ']') {
					// Handle address
					variable, exists := vars[j[1:len(j)-1]]
					if !exists {
						return nil, nil, fmt.Errorf("address %s not declared", j[1:len(j)-1])
					}
					op[i+1] = variable

				} else if num, exists := c.RegNames[j]; exists {
					// Handle register alias
					op[i+1] = num

				} else if len(j) > 0 {
					// Handle int TODO: negatives
					n, err := strconv.Atoi(j)
					if err != nil {
						return nil, nil, err
					}
					op[i+1] = uint16(n)
				}
			}

			ops = append(ops, op)
		}
	}

	return mem, ops, nil
}
