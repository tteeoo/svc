package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	for i := 0; i < c.VGA.TextHeight; i++ {
		for j := 0; j < c.VGA.TextWidth; j++ {
			c.Mem.Mem[uint16(j+(i*c.VGA.TextHeight))] = 0x0f41
			if j == c.VGA.TextWidth-1 {
				c.VGA.TextDraw()
			}
		}
	}
}
