package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	c.Op(c.GetOp("nop"), []uint16{})
	c.Op(c.GetOp("cpl"), []uint16{0, 0xff})
	c.Op(c.GetOp("cop"), []uint16{1, 0})
	c.Op(c.GetOp("str"), []uint16{2, 1})
	c.Op(c.GetOp("ldr"), []uint16{3, 2})
	c.Op(c.GetOp("cpl"), []uint16{6, 0x41})
	c.Op(c.GetOp("str"), []uint16{3, 6})
	fmt.Println(c.GetMem())
}
