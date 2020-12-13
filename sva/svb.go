package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// svb represents a simple virtual binary file.
type svb struct {
	mem map[uint16]uint16
	ops [][]uint16
}

// bytes converts an svb struct to its byte representation.
func (s svb) bytes() ([]byte, error) {
	_, exists := s.mem[0xffff]
	if exists {
		return nil, errors.New("cannot initialize memory address 0xffff in svb file")
	}
	// Calculate length
	l := (len(s.mem) * 2) + 2
	for _, i := range s.ops {
		l += len(i)
	}
	// Create uint16 slice
	b := make([]uint16, l)
	i := 0
	for k, v := range s.mem {
		b[i] = k
		b[i+1] = v
		i += 2
	}
	// Seperator
	b[i] = 0xffff
	b[i+1] = 0xffff
	i += 2
	for _, o := range s.ops {
		for _, p := range o {
			b[i] = p
			i++
		}
	}
	// Convert uint16 to binary
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, b)
	return append([]byte("SVCB"), buf.Bytes()...), nil
}
