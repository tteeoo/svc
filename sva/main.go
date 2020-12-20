package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Get input file
	if len(os.Args) < 2 {
		fmt.Printf("run like this: %s <input file> [-o <output file>]\n", os.Args[0])
		os.Exit(1)
	}
	inputFile := os.Args[1]

	// Get output file
	outputFile := "./out.svb"
	if len(os.Args) == 4 {
		if os.Args[2] == "-o" {
			outputFile = os.Args[3]
		}
	}

	// Read input file
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("error reading:", err)
		os.Exit(1)
	}

	// Parse input
	svbStruct, err := parse(b)
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
