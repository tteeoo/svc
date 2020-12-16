package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	for i := 0; i < 25*80; i++ {
		c.Mem.Set(uint16(i), 0x0f42)
	}
	c.VGA.TextDraw()
}
