// Package cpu implements virtual CPUs.
package cpu

import (
	"github.com/tteeoo/svc/mem"
)

// CPU is an interface that wraps methods to interact with a basic CPU.
type CPU interface {
	GetOp(string) uint16
	Op(uint16, []uint16) error
	GetMemoryDevice() mem.MemoryDevice
}
