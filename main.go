package main

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
)

func main() {
	var m mem.MemoryDevice = mem.NewGenericMemoryDevice()
	fmt.Println(m.GetTextBuffer())
}
