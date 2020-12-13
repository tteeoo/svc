package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewGenericCPU()
	c.Op(c.GetOp("nop"), []uint16{})
	c.Op(c.GetOp("cpl"), []uint16{0xff, 0})
	c.Op(c.GetOp("cop"), []uint16{0, 1})
	c.Op(c.GetOp("str"), []uint16{1, 2, 2})
	c.Op(c.GetOp("ldr"), []uint16{2, 2, 3})
	c.Op(c.GetOp("cpl"), []uint16{0x41, 6})
	c.Op(c.GetOp("str"), []uint16{2, 3, 6})
	fmt.Println(c.GetMemoryDevice().GetAddressSpace(0))
}
