package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Get input file
	if len(os.Args) < 2 {
		fmt.Printf("run like this: %s <input file> [-o <output file>] [-p]\n", os.Args[0])
		os.Exit(1)
	}

	// Determine if we should output the pre-processed assembly
	writePP := false
	var args []string
	for _, a := range os.Args[1:] {
		if a != "-p" {
			args = append(args, a)
		} else {
			writePP = true
		}

	}
	inputFile := args[0]

	// Get output file
	outputFile := "./out.svb"
	if len(args) == 3 {
		if args[1] == "-o" {
			outputFile = args[2]
		}
	}

	// Read input file
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("error reading:", err)
		os.Exit(1)
	}

	// Pre-process input
	lines, err := preProcess(b, true)
	if err != nil {
		fmt.Println("error pre-processing:", err)
		os.Exit(1)
	}

	// Write pre-processed input
	if writePP {
		ppOut := ""
		for _, i := range lines {
			content := false
			for _, j := range i {
				if j != "" {
					ppOut += j + " "
					content = true
				}
			}
			if content {
				ppOut += "\n"
			}
		}

		err = ioutil.WriteFile(outputFile + ".asm",  []byte(ppOut), 0644)
		if err != nil {
			fmt.Println("error writing binary:", err)
			os.Exit(1)
		}
	}

	// Parse input
	svbStruct, err := parse(lines)
	if err != nil {
		fmt.Println("error parsing:", err)
		os.Exit(1)
	}

	// fmt.Printf("%+v\n", svbStruct)

	// Parse binary
	out, err := svbStruct.Bytes()
	if err != nil {
		fmt.Println("error parsing binary:", err)
		os.Exit(1)
	}

	// Write binary
	err = ioutil.WriteFile(outputFile, out, 0644)
	if err != nil {
		fmt.Println("error writing binary:", err)
		os.Exit(1)
	}
}
