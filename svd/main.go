package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/mem"
	"github.com/tteeoo/svc/svb"
	"github.com/tteeoo/svc/vga"
	"io/ioutil"
	"os"
)

func main() {

	// Get program file
	if len(os.Args) < 2 {
		fmt.Printf("run like this: %s <svb file>\n", os.Args[0])
		os.Exit(1)
	}
	programFile := os.Args[1]

	// Open program
	b, err := ioutil.ReadFile(programFile)
	if err != nil {
		fmt.Println("error reading program file:", err)
		os.Exit(1)
	}

	m := mem.NewRAM(mem.AddressSpace{}, 80, 25)
	v := vga.NewVGA(m)
	c := cpu.NewCPU(m, v)

	// Parse program
	fmt.Println("simple virtual debugger: alpha")
	fmt.Printf("parsing file: [%s]... ", os.Args[1])
	program, err := svb.ParseBinary(c, b)
	if err != nil {
		fmt.Println("error parsing program file:", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	// Load program into memory
	m.Mem = program.GetProgramMem()

	// Calculate heap offset
	m.HeapOffset += program.Size()

	// Start repl
	repl(c, program.MainAddress)
}
