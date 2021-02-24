// Package vga implements video devices which can show output.
package vga

import (
	"fmt"
	"github.com/tteeoo/svc/mem"
	"strings"
)

// VGA represents a video device.
type VGA struct {
	// Mem is the memory device to read from.
	Mem *mem.RAM
	// LastBuffer stores the last rendered buffer split into lines.
	LastBuffer []string
}

// NewVGA returns a pointer to a newly initialized VGA.
func NewVGA(m *mem.RAM) *VGA {
	return &VGA{
		LastBuffer: []string{},
		Mem:        m,
	}
}

// TextDraw reads from memory and prints out the text-buffer.
func (v *VGA) TextDraw() {
	// Initialize
	tb := make([][][2]byte, v.Mem.VGAHeight)
	for i := range tb {
		tb[i] = make([][2]byte, v.Mem.VGAWidth)
	}
	// Populate
	a := uint16(0)
	for i := 0; i < v.Mem.VGAHeight; i++ {
		for j := 0; j < v.Mem.VGAWidth; j++ {
			b := v.Mem.Get(a)
			tb[i][j] = [2]byte{
				byte(b >> 8),
				byte((b << 8) >> 8),
			}
			a++
		}
	}
	// Print
	out := ""
	for _, i := range tb {
		for _, j := range i {
			attr := [2]int{
				int(j[0] >> 4),
				int((j[0] << 4) >> 4),
			}
			// Translate 16 VGA colors to ANSI colors
			switch attr[0] {
			case 0:
				out += "\033[40m"
			case 1:
				out += "\033[44m"
			case 2:
				out += "\033[42m"
			case 3:
				out += "\033[46m"
			case 4:
				out += "\033[41m"
			case 5:
				out += "\033[45m"
			case 6:
				out += "\033[43m"
			case 7:
				out += "\033[47m"
			case 8:
				out += "\033[100m"
			case 9:
				out += "\033[104m"
			case 10:
				out += "\033[102m"
			case 11:
				out += "\033[106m"
			case 12:
				out += "\033[101m"
			case 13:
				out += "\033[105m"
			case 14:
				out += "\033[103m"
			case 15:
				out += "\033[107m"
			}
			switch attr[1] {
			case 0:
				out += "\033[30m"
			case 1:
				out += "\033[34m"
			case 2:
				out += "\033[32m"
			case 3:
				out += "\033[36m"
			case 4:
				out += "\033[31m"
			case 5:
				out += "\033[35m"
			case 6:
				out += "\033[33m"
			case 7:
				out += "\033[37m"
			case 8:
				out += "\033[90m"
			case 9:
				out += "\033[94m"
			case 10:
				out += "\033[92m"
			case 11:
				out += "\033[96m"
			case 12:
				out += "\033[91m"
			case 13:
				out += "\033[95m"
			case 14:
				out += "\033[93m"
			case 15:
				out += "\033[97m"
			}
			out += string(j[1]) + "\033[0m"
		}
		out += "\n"
	}
	current := strings.Split(out, "\n")
	if len(v.LastBuffer) != len(current) {
		print("\033[2J\033[H" + out)
		v.LastBuffer = current
		return
	}
	realOut := "\033[H"
	diffLines := []int{}
	for i, j := range current {
		for k := range j {
			if current[i][k] != v.LastBuffer[i][k] {
				diffLines = append(diffLines, i)
				break
			}
		}
	}
	v.LastBuffer = current
	line := 0
	for _, i := range diffLines {
		cursor := (i - line)
		if cursor > 0 {
			realOut += fmt.Sprintf("\033[%dB", cursor)
		} else if cursor < 0 {
			realOut += fmt.Sprintf("\033[%dA", -cursor)
		}
		realOut += "\033[2K" + current[i] + "\n\033[0m"
		line = i + 1
	}
	final := ""
	if v.Mem.VGAHeight-line > 0 {
		final = fmt.Sprintf("\033[%dB", v.Mem.VGAHeight-line)
	}
	print(realOut + final)
}
