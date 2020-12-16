package main

import (
	"github.com/tteeoo/svc/cpu"
)

func main() {
	c := cpu.NewCPU()
	c.Op(c.GetOp("nop"), []uint16{})
	c.Op(c.GetOp("cpl"), []uint16{c.RegNames["ac"], 4})
	c.Op(c.GetOp("cpl"), []uint16{0, 4})
	c.Op(c.GetOp("add"), []uint16{0})
}
