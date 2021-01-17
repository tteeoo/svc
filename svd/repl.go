package main

import (
	"fmt"
	"github.com/tteeoo/svc/cpu"
	"github.com/tteeoo/svc/dat"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

var done bool

func run(c *cpu.CPU) bool {
	pc := c.Regs[dat.RegNamesToNum["pc"]]

	// Exit if pc is the last address
	if pc == 0xffff {
		fmt.Println("program counter is ffff, execution stopped")
		done = true
		return true
	}

	// Get instruction
	op := c.Mem.Get(pc)
	name := dat.OpCodeToName[op]
	size := dat.OpNameToSize[name]
	operands := make([]uint16, size)
	for i := uint16(0); i < size; i++ {
		operands[i] = c.Mem.Get(pc + (1 + i))
	}
	fmt.Println(
		color(fmt.Sprintf("%x:", pc), "32;1"),
		color(fmt.Sprintf("%s(%x)", name, op), "31;1"),
		color(fmt.Sprintf("%x", operands), "31;1"),
	)

	// Set the lc register
	c.Regs[dat.RegNamesToNum["lc"]] = pc

	// Increase program counter
	c.Regs[dat.RegNamesToNum["pc"]] += uint16(1 + size)

	// Execute instruction
	if op == dat.OpNameToCode["vga"] {
		fmt.Println(color("text drawn", "35;1"))
	} else {
		c.Op(op, operands)
	}

	return false
}

func repl(c *cpu.CPU, address uint16) {

	// Put command-line args into heap
	var l uint16
	if len(os.Args) > 2 {
		i := c.Mem.HeapOffset
		for _, str := range os.Args[2:] {
			for _, char := range str {
				c.Mem.Set(i, uint16(char))
				i++
			}
			c.Mem.Set(i, 0)
			i++
		}
		fmt.Println(color(fmt.Sprintf("argument(s) loaded into heap: %s", os.Args[2:]), "33;1"))
		l = uint16(i - c.Mem.HeapOffset)
	}

	// Load heap information
	c.Mem.Set(0xffff, l)
	c.Mem.Set(0xfffe, c.Mem.HeapOffset)

	// Push exit address onto stack
	sp := dat.RegNamesToNum["sp"]
	c.Mem.Set(c.Regs[sp], 0xffff)
	fmt.Println(color("pushed ffff onto stack", "36;1"))

	// Set the program counter
	c.Regs[dat.RegNamesToNum["pc"]] = address
	fmt.Println(color(fmt.Sprintf("program counter set to %x", address), "32;1"))
	fmt.Println("run `h` for help")

	// Enter the execution loop
	var cycles int
	rl, _ := readline.New("> ")
	for {
		// Read input
		input, err := rl.Readline()
		if err != nil {
			panic(err)
		}
		command := strings.Split(strings.TrimSpace(input), " ")

		// Handle command
		switch command[0] {
		// Empty input
		case "":
			if done {
				continue
			}
			run(c)
			cycles++
		// CPU
		case "c":
			// Print registers
			if len(command) == 1 {
				// Sort keys
				keys := make([]string, len(dat.RegNamesToNum))
				for k, v := range dat.RegNamesToNum {
					keys[v] = k
				}
				for _, k := range keys {
					v := dat.RegNamesToNum[k]
					var ansic string
					if k == "pc" || k == "lc" {
						ansic = "32;1"
					} else if k == "sp" {
						ansic = "36;1"
					} else {
						ansic = "34;1"
					}
					fmt.Println(
						color(fmt.Sprintf("%s(%x):", k, v), ansic),
						color(fmt.Sprintf("%x", c.Regs[v]), ansic),
					)
				}
			} else if len(command) == 3 {
				// Change value
				idx, valid := dat.RegNamesToNum[command[1]]
				if !valid {
					fmt.Println("invalid register")
					continue
				}
				value, err := parseHex(command[2])
				if err != nil {
					fmt.Println("invalid value")
					continue
				}
				c.Regs[idx] = value
				fmt.Printf("set register %s to %s\n", command[1], command[2])
			} else {
				fmt.Println("invalid command")
			}
		// Memory
		case "m":
			// Print sections
			if len(command) == 1 {
				fmt.Println(color(fmt.Sprintf("%s: %x-%x", "text", 0, c.Mem.StackMin-1), "35;1"))
				fmt.Println(color(fmt.Sprintf("%s: %x-%x", "stak", c.Mem.StackMin, c.Mem.StackMax), "36;1"))
				fmt.Println(color(fmt.Sprintf("%s: %x-%x, main: %x", "prog", c.Mem.ProgramOffset, c.Mem.HeapOffset-1, address), "32;1"))
				fmt.Println(color(fmt.Sprintf("%s: %x-%x", "heap", c.Mem.HeapOffset, 0xffff), "33;1"))
			} else if len(command) == 2 {
				// Print memory
				memRange := strings.Split(command[1], "-")
				if len(memRange) == 1 {
					value, err := parseHex(memRange[0])
					if err != nil {
						fmt.Println("invalid address")
						continue
					}
					fmt.Println(
						color(
							fmt.Sprintf("%s: %x", command[1], c.Mem.Get(value)),
							sectionToColor(addressToSection(c, value)),
						),
					)
				} else if len(memRange) == 2 {
					// Print memory range
					start, err := parseHex(memRange[0])
					if err != nil {
						fmt.Println("invalid range")
						continue
					}
					end, err := parseHex(memRange[1])
					if err != nil {
						fmt.Println("invalid range")
						continue
					}
					for i := start; i < end+1; i++ {
						if (i-start)%8 == 0 {
							fmt.Print(
								color(
									fmt.Sprintf("%0*x: ", 4, i),
									sectionToColor(addressToSection(c, i)),
								),
							)
						}
						fmt.Print(
							color(
								fmt.Sprintf("%0*x ", 4, c.Mem.Get(i)),
								sectionToColor(addressToSection(c, i)),
							),
						)
						if (i-start+1)%8 == 0 {
							fmt.Print("\n")
						}
					}
					fmt.Print("\n")
				} else {
					fmt.Println("invalid command")
				}
			} else if len(command) == 3 {
				// Set memory
				key, err := parseHex(command[1])
				if err != nil {
					fmt.Println("invalid address")
					continue
				}
				value, err := parseHex(command[2])
				if err != nil {
					fmt.Println("invalid value")
					continue
				}
				c.Mem.Set(key, value)
				fmt.Println(
					color(
						fmt.Sprintf("set memory address %s to %s", command[1], command[2]),
						sectionToColor(addressToSection(c, key)),
					),
				)
			} else {
				fmt.Println("invalid command")
			}
		case "n":
			fmt.Println(color(fmt.Sprintf("%d", cycles), "34;1"))
		case "h", "?":
			fmt.Println("h      print this help message")
			fmt.Println("n      print the number of clock cycles executed")
			fmt.Println("<num>  execute a number of instructions")
			fmt.Println("c                print cpu registers")
			fmt.Println("c <reg> <value>  set cpu register")
			fmt.Println("m                 print the sections of memory")
			fmt.Println("m <addr>          print memory address")
			fmt.Println("m <addr>-<addr>   print range of memory")
			fmt.Println("m <addr> <value>  set memory address")
		default:
			// Try number
			if len(command) == 1 {
				num, err := strconv.Atoi(command[0])
				if err == nil {
					if done {
						fmt.Println("execution has stopped")
						continue
					}
					for i := 0; i < num; i++ {
						cycles++
						if run(c) {
							break
						}
					}
				} else {
					fmt.Println("invalid command")
				}
			} else {
				fmt.Println("invalid command")
			}
		}
	}
}
