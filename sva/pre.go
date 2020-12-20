package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"path"
	"os"
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
		for _, str := range badSplitLine {
			if str != "" {
				splitLine = append(splitLine, str)
			}
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
