package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	a := uint16(0)
	for i := 0; i < c.VGA.TextHeight; i++ {
		for j := 0; j < c.VGA.TextWidth; j++ {
			c.Mem.Set(a, (uint16(j)<<8)+0x41)
			if j%c.VGA.TextWidth == 0 {
				c.VGA.TextDraw()
			}
			a++
		}
	}
	for i := c.VGA.TextHeight - 1; i >= 0; i-- {
		for j := 0; j < c.VGA.TextWidth; j++ {
			c.Mem.Set(a, (((uint16(j)<<8)<<4)+((uint16(j)<<8)>>4)<<4)+0x42)
			if j%c.VGA.TextWidth == 0 {
				c.VGA.TextDraw()
			}
			a--
		}
	}
	for x := 0; x < c.VGA.TextHeight; x++ {
		for i, j := range "Hello, world!" {
			if x%2 != 0 {
				c.Mem.Set(uint16((80*x)+i), 0x0f00+uint16(j))
			}
		}
		c.VGA.TextDraw()
	}
	for x := c.VGA.TextHeight - 1; x >= 0; x-- {
		for i, j := range "Hello, world!" {
			if x%2 == 0 {
				c.Mem.Set(uint16((80*x)+i+(c.VGA.TextWidth-13)), 0x0f00+uint16(j))
			}
		}
		c.VGA.TextDraw()
	}
}
