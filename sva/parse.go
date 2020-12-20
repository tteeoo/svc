package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tteeoo/svc/dat"
	"github.com/tteeoo/svc/svb"
	"strconv"
	"strings"
)

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

// parseNum will take a number-representing string and parse it to a
// two's complement uint16.
func parseNum(s string) (uint16, error) {

	// Negative number
	if s[0] == '-' {
		n, err := strconv.Atoi(s[1:])
		if n > 65536 {
			return 0, fmt.Errorf("int value \"%s\" is too large", s)
		}
		if err != nil {
			return 0, err
		}
		return ^uint16(n) + 1, nil
	}

	// Positive number
	n, err := strconv.Atoi(s)
	if n > 65536 {
		return 0, fmt.Errorf("int value \"%s\" is too large", s)
	}
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

// parse will parse a pre-processed input file into an SVB struct.
func parse(b []byte) (svb.SVB, error) {

	vars := make(map[string]uint16)
	subs := make(map[string]uint16)
	address := uint16(0)
	constants := []svb.Constant{}
	currentSub := svb.Subroutine{Size: -1}
	binary := svb.SVB{}

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
					// Create constants for each char
					vars[splitLine[0]] = address
					for _, char := range splitLine[2][1 : len(splitLine[2])-1] {
						constants = append(constants, svb.Constant{
							Address: address,
							Value:   uint16(char),
						})
						address++
					}
					continue

				} else if splitLine[2][1] == 'x' {
					// Handle a hex value
					val, err := parseHex(splitLine[2][2:])
					if err != nil {
						return svb.SVB{}, err
					}
					// Create constant
					vars[splitLine[0]] = address
					constants = append(constants, svb.Constant{
						Address: address,
						Value:   val,
					})
					address++
					continue
				}
			}

			// Handle an int
			i, err := parseNum(splitLine[2])
			if err != nil {
				return svb.SVB{}, err
			}
			vars[splitLine[0]] = address
			constants = append(constants, svb.Constant{
				Address: address,
				Value:   uint16(i),
			})
			address++

		} else if len(splitLine) == 1 && len(splitLine[0]) > 1 && splitLine[0][len(splitLine[0])-1] == ':' {
			// Handle subroutine definition
			name := splitLine[0][:len(splitLine[0])-1]
			if _, exists := subs[name]; exists {
				return svb.SVB{}, fmt.Errorf("subroutine \"%s\" already exists", name)
			}
			if currentSub.Size != -1 {
				binary.Subroutines = append(binary.Subroutines, currentSub)
			}

			subs[name] = address
			currentSub = svb.Subroutine{
				Name:    name,
				Address: address,
			}
			address++

		} else if len(splitLine) > 0 {

			// Handle instruction
			code, exists := dat.OpNameToCode[splitLine[0]]
			if !exists {
				return svb.SVB{}, fmt.Errorf("instruction \"%s\" does not exist", splitLine[0])
			}
			operands := make([]uint16, len(splitLine)-1)

			for i, j := range splitLine[1:] {
				if (len(j) > 2) && (j[1] == 'x') {
					// Handle hex number
					val, err := parseHex(j[2:])
					if err != nil {
						return svb.SVB{}, err
					}
					operands[i] = val

				} else if (len(j) > 2) && (j[0] == '[') && (j[len(j)-1] == ']') {
					// Handle constant reference
					variable, exists := vars[j[1:len(j)-1]]
					if !exists {
						return svb.SVB{}, fmt.Errorf("constant \"%s\" not declared", j[1:len(j)-1])
					}
					operands[i] = variable

				} else if (len(j) > 2) && (j[0] == '{') && (j[len(j)-1] == '}') {
					// Handle constant reference
					subAddr, exists := subs[j[1:len(j)-1]]
					if !exists {
						return svb.SVB{}, fmt.Errorf("subroutine \"%s\" not declared", j[1:len(j)-1])
					}
					operands[i] = subAddr

				} else if num, exists := dat.RegNamesToNum[j]; exists {
					// Handle register alias
					operands[i] = num

				} else if len(j) > 0 {
					// Handle an int
					num, err := parseNum(j)
					if err != nil {
						return svb.SVB{}, err
					}
					operands[i] = num
				}
			}

			// Check that the right number of operands are provided
			size := dat.OpNameToSize[splitLine[0]]
			if len(operands) != size {
				return svb.SVB{}, fmt.Errorf("operation %s expected %d operands, but received %d", splitLine, size, len(operands))
			}

			// Check to make sure instruction is in a defined subroutine
			if currentSub.Size == -1 {
				return svb.SVB{}, fmt.Errorf("instruction %s used outside of a subroutine", splitLine)
			}

			currentSub.Instructions = append(currentSub.Instructions, svb.Instruction{
				Name:     splitLine[0],
				Opcode:   code,
				Operands: operands,
				Size:     size + 1,
			})
			currentSub.Size += size + 1
		}
	}

	if currentSub.Size != -1 {
		binary.Subroutines = append(binary.Subroutines, currentSub)
	}
	binary.Constants = constants

	// TODO: fix overlapping addresses and set main subroutine address

	return binary, nil
}
