package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewGenericCPU()
	c.Op(0x00, []uint16{})
}
