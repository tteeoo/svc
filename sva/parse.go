package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/dat"
	"github.com/tteeoo/svc/svb"
	"github.com/tteeoo/svc/util"
	"strconv"
)

// parseNum will take a number-representing string and parse it to a
//   two's complement uint16.
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
func parse(c *cpu.CPU, lines [][]string) (svb.SVB, error) {

	vars := make(map[string]uint16)
	subs := make(map[string]uint16)
	labelIndices := make(map[string][][3]int)
	labelAddresses := make(map[string]uint16)
	address := c.Mem.ProgramOffset
	constants := []svb.Constant{}
	currentSub := svb.Subroutine{}
	binary := svb.SVB{}

	// Iterate lines
	for _, splitLine := range lines {

		// Handle constants
		if (len(splitLine) == 3) && (splitLine[1] == "=") {
			if currentSub.Name != "" {
				return svb.SVB{},
					fmt.Errorf("you cannot define a constant inside of a subroutine (\"%s\" is in \"%s\")",
						splitLine,
						currentSub.Name,
					)
			}
			if _, exists := vars[splitLine[0]]; exists {
				return svb.SVB{}, fmt.Errorf("constant \"%s\" defined more than once", splitLine[0])
			}
			if len(splitLine[2]) > 2 {

				// Handle a string
				if (splitLine[2][0] == byte('"')) && (splitLine[2][len(splitLine[2])-1] == byte('"')) {
					// Create constants for each char
					vars[splitLine[0]] = address
					for _, char := range splitLine[2][1 : len(splitLine[2])-1] {
						constants = append(constants, svb.Constant{
							Name:    splitLine[0],
							Address: address,
							Value:   uint16(char),
						})
						address++
					}
					constants = append(constants, svb.Constant{
						Name:    splitLine[0],
						Address: address,
						Value:   uint16(0),
					})
					address++
					continue

				} else if splitLine[2][1] == 'x' {
					// Handle a hex value
					val, err := util.ParseHex(splitLine[2][2:])
					if err != nil {
						return svb.SVB{}, err
					}
					// Create constant
					vars[splitLine[0]] = address
					constants = append(constants, svb.Constant{
						Name:    splitLine[0],
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
				Name:    splitLine[0],
				Address: address,
				Value:   uint16(i),
			})
			address++

		} else if len(splitLine) == 1 && len(splitLine[0]) > 1 && splitLine[0][0] == '&' {
			// Handle label definition
			name := splitLine[0][1:]
			if _, exists := labelAddresses[name]; exists {
				return svb.SVB{}, fmt.Errorf("label \"%s\" defined more than once", name)
			}
			if currentSub.Name == "" {
				return svb.SVB{}, fmt.Errorf("label \"%s\" defined outside of a subroutine", name)
			}

			labelAddresses[name] = address + uint16(currentSub.Size())

		} else if len(splitLine) == 1 && len(splitLine[0]) > 1 && splitLine[0][len(splitLine[0])-1] == ':' {
			// Handle subroutine definition
			name := splitLine[0][:len(splitLine[0])-1]
			if _, exists := subs[name]; exists {
				return svb.SVB{}, fmt.Errorf("subroutine \"%s\" defined more than once", name)
			}
			if currentSub.Name != "" {
				binary.Subroutines = append(binary.Subroutines, currentSub)
				address += uint16(currentSub.Size())
			}

			subs[name] = address
			currentSub = svb.Subroutine{
				Name:    name,
				Address: address,
			}

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
					val, err := util.ParseHex(j[2:])
					if err != nil {
						return svb.SVB{}, err
					}
					operands[i] = val

				} else if (len(j) > 1) && (j[0] == '&') {
					// Handle label reference
					name := j[1:]
					labelIndices[name] = append(labelIndices[name], [3]int{
						len(binary.Subroutines),
						len(currentSub.Instructions),
						i,
					})
					operands[i] = 0

				} else if (len(j) > 2) && (j[0] == '[') && (j[len(j)-1] == ']') {
					// Handle constant reference
					variable, exists := vars[j[1:len(j)-1]]
					if !exists {
						return svb.SVB{}, fmt.Errorf("constant \"%s\" not declared", j[1:len(j)-1])
					}
					operands[i] = variable

				} else if (len(j) > 2) && (j[0] == '{') && (j[len(j)-1] == '}') {
					// Handle subroutine reference
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
			if len(operands) != size+dat.OpNameToPacked[splitLine[0]] {
				return svb.SVB{},
					fmt.Errorf("operation \"%s\" expected %d operands, but received %d",
						splitLine,
						size,
						len(operands),
					)
			}

			// Check to make sure instruction is in a defined subroutine
			if currentSub.Name == "" {
				return svb.SVB{}, fmt.Errorf("instruction \"%s\" used outside of a subroutine", splitLine)
			}

			currentSub.Instructions = append(currentSub.Instructions, svb.Instruction{
				Name:     splitLine[0],
				Opcode:   code,
				Operands: operands,
			})
		}
	}

	// Handle main routine
	if currentSub.Name != "main" {
		return svb.SVB{}, fmt.Errorf("the last subroutine \"%s\", is not named \"main\"", currentSub.Name)
	}
	binary.MainAddress = currentSub.Address

	binary.Subroutines = append(binary.Subroutines, currentSub)
	binary.Constants = constants

	// Set label addresses
	for k, v := range labelIndices {
		for _, ref := range v {
			binary.Subroutines[ref[0]].Instructions[ref[1]].Operands[ref[2]] = labelAddresses[k]
		}
	}

	return binary, nil
}
