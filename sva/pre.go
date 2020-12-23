package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Handle ex register expansion "ldr (0) aa" -> "cpl ex 0", "ldr ex aa"
func registerExpansion(splitLine []string, idx int) ([][]string, error) {
	first := []string{"cpl", "ex"}

	if splitLine[idx][0] != '(' || splitLine[idx][len(splitLine[idx])-1] != ')' {
		return [][]string{}, fmt.Errorf("invalid register expansion: %s", splitLine)
	}
	first = append(first, splitLine[idx][1:len(splitLine[idx])-1])
	splitLine[idx] = "ex"

	return [][]string{first, splitLine}, nil
}

// Handle double operator "aaa, bbb xx yy" -> "aaa xx yy", "bbb xx yy"
func expandOperation(splitLine []string) [][]string {
	first := []string{splitLine[0][:len(splitLine[0])-1]}
	second := []string{splitLine[1]}

	if len(splitLine) > 2 {
		for _, and := range splitLine[2:] {
			first = append(first, and)
			second = append(second, and)
		}
	}

	return [][]string{first, second}
}

// Handle double operands "aaa bb xx, yy zz" -> "aaa bb xx", "aaa yy zz"
func expandOperand(splitLine []string) ([][]string, error) {
	first := []string{splitLine[0]}
	second := []string{splitLine[0]}

	if len(splitLine) < 3 {
		return [][]string{}, fmt.Errorf("invalid instruction expansion: %s", splitLine)
	}
	onFirst := true
	for _, and := range splitLine[1:] {
		if !onFirst {
			second = append(second, and)
			continue
		}
		if and[len(and)-1] != ',' {
			first = append(first, and)
			continue
		}
		first = append(first, and[:len(and)-1])
		onFirst = false
	}

	return [][]string{first, second}, nil
}

// returns an expansion if exists, else returns the original instruction
func detectExpansion(splitLine []string) ([][]string, error) {
	for i, s := range splitLine {
		for j, c := range s {
			if c == '=' {
				return [][]string{splitLine}, nil
			}
			if c == '(' && j == 0 && i != 0 && len(s) > 1 {
				return registerExpansion(splitLine, i)
			}
			if c == ',' {
				if len(splitLine) < 2 {
					return [][]string{}, fmt.Errorf("invalid instruction expansion: %s", splitLine)
				}
				if i == 0 {
					return expandOperation(splitLine), nil
				}
				return expandOperand(splitLine)
			}
		}
	}
	return [][]string{splitLine}, nil
}

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

		// Detect expansions
		expandedLines, err := detectExpansion(splitLine)
		if err != nil {
			return [][]string{}, err
		}
		if len(expandedLines) != 1 {
			for _, l := range expandedLines {
				lines = append(lines, l)
			}
			continue
		}

		// Normal instruction
		if len(splitLine) != 2 || splitLine[0] != "." {
			lines = append(lines, splitLine)
			continue
		}

		// Handle file sourcing
		if !allowSource {
			return [][]string{}, fmt.Errorf("cannot recursively source files (attempting to source %s)", splitLine[1])
		}
		var fb []byte

		// Handle absolute path
		if path.IsAbs(splitLine[1]) {
			fb, err = ioutil.ReadFile(splitLine[1])
			if err != nil {
				return [][]string{}, err
			}
		} else {
			// Handle relative path
			wd, err := os.Getwd()
			if err != nil {
				return [][]string{}, err
			}
			fb, err = ioutil.ReadFile(path.Join(wd, splitLine[1]))
			if err != nil {
				return [][]string{}, err
			}
		}

		// Append lines from file
		flines, err := preProcess(fb, false)
		if err != nil {
			return [][]string{}, err
		}
		lines = append(lines, flines...)
	}

	return lines, nil
}
