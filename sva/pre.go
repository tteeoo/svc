package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// preProcess will preProcess an assembly file.
// It will remove comments and expand file sources.
func preProcess(b []byte, allowSource bool) ([][]string, error) {

	var lines [][]string
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
		inString := false
		for _, str := range badSplitLine {
			if str != "" {
				// Handle spaces in strings
				if inString {
					splitLine[len(splitLine)-1] += " " + str
					continue
				}
				if str[0] == '"' {
					inString = true
				}
				splitLine = append(splitLine, str)
			}
		}

		// Detect instruction expansions
		doubleOp := false
		doubleAnd := false
		isConstant := false
		for i, s := range splitLine {
			for _, c := range s {
				if c == '=' {
					isConstant = true
				} else if c == ',' {
					if len(splitLine) < 2 || doubleAnd || doubleOp {
						return [][]string{}, fmt.Errorf("invalid instruction expansion: %s\n", splitLine)
					}
					if !isConstant {
						if i == 0 {
							doubleOp = true
						} else {
							doubleAnd = true
						}
					}
				}
			}
		}

		// Handle double operator (aaa, bbb xx yy) -> (aaa xx yy), (bbb xx yy)
		if doubleOp {
			first := []string{splitLine[0][:len(splitLine[0])-1]}
			second := []string{splitLine[1]}

			if len(splitLine) > 2 {
				for _, and := range splitLine[2:] {
					first = append(first, and)
					second = append(second, and)
				}
			}

			lines = append(append(lines, first), second)
			continue

		} else if doubleAnd {
			// Handle double operands (aaa bb xx, yy zz) -> (aaa bb xx), (aaa yy zz)
			first := []string{splitLine[0]}
			second := []string{splitLine[0]}

			if len(splitLine) < 3 {
				return [][]string{}, fmt.Errorf("invalid instruction expansion: %s\n", splitLine)
			}
			onFirst := true
			for _, and := range splitLine[1:] {
				if onFirst {
					if and[len(and)-1] == ',' {
						first = append(first, and[:len(and)-1])
						onFirst = false
					} else {
						first = append(first, and)
					}
				} else {
					second = append(second, and)
				}
			}

			lines = append(append(lines, first), second)
			continue
		}

		// Handle file sourcing
		if len(splitLine) == 2 && splitLine[0] == "." {
			if allowSource {
				if path.IsAbs(splitLine[1]) {

					// Handle absolute path
					fb, err := ioutil.ReadFile(splitLine[1])
					if err != nil {
						return [][]string{}, err
					}

					// Append processed lines
					flines, err := preProcess(fb, false)
					if err != nil {
						return [][]string{}, err
					}
					lines = append(lines, flines...)
				} else {

					// Handle relative path
					wd, err := os.Getwd()
					if err != nil {
						return [][]string{}, err
					}
					fb, err := ioutil.ReadFile(path.Join(wd, splitLine[1]))
					if err != nil {
						return [][]string{}, err
					}

					// Append processed lines
					flines, err := preProcess(fb, false)
					if err != nil {
						return [][]string{}, err
					}
					lines = append(lines, flines...)
				}
			} else {
				return [][]string{},
					fmt.Errorf(
						"cannot recursively source files (attempting to source %s)",
						splitLine[1],
					)
			}
		} else {
			lines = append(lines, splitLine)
		}
	}

	return lines, nil
}
