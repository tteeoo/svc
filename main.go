package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/dat"
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

	// Parse program
	program, err := svb.ParseBinary(b)
	if err != nil {
		fmt.Println("error parsing program file:", err)
		os.Exit(1)
	}

	// Initialize CPU
	space := program.GetProgramMem()
	m := mem.NewRAM(space)
	v := vga.NewVGA(m, dat.VGAOffset, dat.VGAHeight, dat.VGAWidth)
	c := cpu.NewCPU(m, v)

	// Calculate heap offset
	dat.HeapOffset += program.Size()

	// Run!
	c.Run(program.MainAddress)
}
