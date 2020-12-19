package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tteeoo/svc/dat"
	"github.com/tteeoo/svc/svb"
	"strconv"
	"strings"
)

// parse will parse a pre-processed input file into an SVB struct.
func parse(b []byte) (svb.SVB, error) {

	var ops [][]uint16
	var address uint16 = 0x00
	constants := []svb.Constant{}
	insructions := []svb.Instruction{}

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

		// Handle constants
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
						return svb.SVB{}, err
					}
					mem[address] = binary.BigEndian.Uint16(b)
					address++
					continue
				}
			}

			// Handle an int
			i, err := strconv.Atoi(splitLine[2])
			if err != nil {
				return svb.SVB{}, err
			}
			vars[splitLine[0]] = address
			mem[address] = uint16(i)
			address++

		} else if len(splitLine) > 0 {
			// Handle instruction
			op := make([]uint16, len(splitLine))
			code, exists := dat.OpNameToCode[splitLine[0]]
			if !exists {
				return svb.SVB{}, fmt.Errorf("instruction \"%s\" does not exist", splitLine[0])
			}
			op[0] = code

			for i, j := range splitLine[1:] {
				if (len(j) > 2) && (j[1] == 'x') {
					// Handle hex number
					val := j[2:]
					if len(val) == 2 {
						val = "00" + val
					}
					b, err := hex.DecodeString(val)
					if err != nil {
						return svb.SVB{}, err
					}
					op[i+1] = binary.BigEndian.Uint16(b)

				} else if (len(j) > 2) && (j[0] == '[') && (j[len(j)-1] == ']') {
					// Handle address
					variable, exists := vars[j[1:len(j)-1]]
					if !exists {
						return svb.SVB{}, fmt.Errorf("address %s not declared", j[1:len(j)-1])
					}
					op[i+1] = variable

				} else if num, exists := dat.RegNamesToNum[j]; exists {
					// Handle register alias
					op[i+1] = num

				} else if len(j) > 0 {
					// Negative number
					if j[0] == '-' && len(j) > 1 {
						n, err := strconv.Atoi(j[1:])
						if err != nil {
							return svb.SVB{}, err
						}
						op[i+1] = ^uint16(n) + 1
					} else {
						// Positive number
						n, err := strconv.Atoi(j)
						if err != nil {
							return svb.SVB{}, err
						}
						op[i+1] = uint16(n)
					}
				}
			}

			ops = append(ops, op)
		}
	}

	return svbStruct, nil
}
