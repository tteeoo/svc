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

	// Initialize devices
	m := mem.NewRAM(mem.AddressSpace{}, 80, 25)
	v := vga.NewVGA(m)
	c := cpu.NewCPU(m, v)

	// Load program into memory
	mainAddress := uint16(0)
	programSize := uint16(0)
	m.Mem, mainAddress, programSize = svb.LoadProgram(c, b)

	// Calculate heap offset
	m.HeapOffset += programSize

	// Run!
	c.Run(mainAddress)
}
