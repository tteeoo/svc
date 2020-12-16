package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	a := uint16(0)
	for i := 0; i < c.VGA.TextHeight; i++ {
		for j := 0; j < c.VGA.TextWidth; j++ {
			c.Mem.Set(a, (uint16(j) << 8) + 0x41)
			a++
		}
	}
	for i, j := range "Hello, world!" {
		c.Mem.Set(uint16(i), 0x0f00 + uint16(j))
	}
	c.VGA.TextDraw()
}
