# Simple Virtual Computer

Currently under development; this is a work in progress.

A 16-bit virtual machine written in Go.

It uses word-based memory, where each memory address maps to a word-length value.
This means that while there are only 65,536 memory addresses,
it technically has 128K of memory since each address points to a 16-bit value instead of a byte.
The reasoning for this is simple: simplicity. Most modern computer systems use bytes because they can be more flexible and efficient.
Those are not the goals of this project.

It implements a "VGA text mode" that reads the contents of memory, using 2,000 contiguous words (which is interpreted as a 80x25 character display).
It translates the encoded VGA text colors into ANSI escape codes and prints the colorized ASCII text.

This repository contains the virtual machine, an assembler to compile programs for it, and a debugger for those programs.

## Instruction Set

The instruction set resembles a RISC (Reduced Instruction Set Computer) architecture.

In the opcode, `r` represents the number of a CPU register that is packed into the word which contains the opcode to save memory.

| Opcode   | Name  | Operands                               | Description                                                                                                                                               |
| -------- | ----- | -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `0x0000` | `nop` |                                        | Does nothing.                                                                                                                                             |
| `0x01rr` | `cop` | `reg to copy to` `reg to copy from`    | Copies the value from one register to another.                                                                                                            |
| `0x020r` | `cpl` | `reg` `value`                          | Copies a literal value to a register.                                                                                                                     |
| `0x03rr` | `str` | `reg holding addr` `reg holding value` | Stores the value from a register into memory at the address held in another register.                                                                     |
| `0x04rr` | `ldr` | `reg to load to` `reg holding addr`    | Loads the value from memory at the address held in one register into another register.                                                                    |
| `0x050r` | `add` | `reg`                                  | Adds the value held in a register to the accumulator.                                                                                                     |
| `0x060r` | `sub` | `reg`                                  | Subtracts the value held in a register from the accumulator.                                                                                              |
| `0x070r` | `twc` | `reg`                                  | Sets a register to the two's complement of the value it holds.                                                                                            |
| `0x080r` | `inc` | `reg`                                  | Increases the value held in a register by one.                                                                                                            |
| `0x090r` | `dec` | `reg`                                  | Decreases the value held in a register by one.                                                                                                            |
| `0x0A0r` | `mul` | `reg`                                  | Multiplies the accumulator with the value held in a register.                                                                                             |
| `0x0B0r` | `div` | `reg`                                  | Divides the accumulator by the value held in a register, storing the modulus in the `ex` register.                                                        |
| `0x0C0r` | `dcv` | `reg`                                  | Divides the accumulator by the value held in a register, taking into account two's complement negative numbers, storing the modulus in the `ex` register. |
| `0x0D0r` | `xor` | `reg`                                  | Performs the "xor" operation on the accumulator with the value held in a register.                                                                        |
| `0x0E0r` | `and` | `reg`                                  | Performs the bitwise "and" operation on the accumulator with the value held in a register.                                                                |
| `0x0F0r` | `orr` | `reg`                                  | Performs the bitwise "or" operation on the accumulator with the value held in a register.                                                                 |
| `0x100r` | `not` | `reg`                                  | Inverts the value held in a register.                                                                                                                     |
| `0x110r` | `shr` | `reg` `value`                          | Shifts the value held in a register to the right by a specified value.                                                                                    |
| `0x120r` | `shl` | `reg` `value`                          | Shifts the value held in a register to the left by a specified value.                                                                                     |
| `0x1300` | `vga` |                                        | Prints the VGA text buffer to the screen.                                                                                                                 |
| `0x140r` | `psh` | `reg`                                  | Decreases the stack pointer and sets the top value of the stack to the value held in a register.                                                          |
| `0x150r` | `pop` | `reg`                                  | Stores the top value of the stack in a register and increases the stack pointer.                                                                          |
| `0x1600` | `ret` |                                        | Pops the program counter off of the stack.                                                                                                                |
| `0x1700` | `cal` | `addr`                                 | Pushes the program counter onto the stack and sets the program counter to an address.                                                                     |
| `0x18rr` | `cmp` | `reg` `reg`                            | If the values of two registers are the same, the boolean index (`bi` register) is set to `0xffff`, else `0xfffe`.                                         |
| `0x1900` | `cle` | `addr`                                 | Equivalent to `cal`, but only executes if the `bi` register is set to `0xffff`.                                                                           |
| `0x1A00` | `cln` | `addr`                                 | Equivalent to `cal`, but only executes if the `bi` register is set to `0xfffe`.                                                                           |
| `0x1B00` | `gto` | `addr`                                 | Sets the program counter to an address.                                                                                                                   |
| `0x1C00` | `gte` | `addr`                                 | Equivalent to `gto`, but only executes if the `bi` register is set to `0xffff`.                                                                           |
| `0x1D00` | `gtn` | `addr`                                 | Equivalent to `gto`, but only executes if the `bi` register is set to `0xfffe`.                                                                           |
| `0x1E0r` | `cml` | `reg` `value`                          | Equivalent to `cmp`, but the second operand is a literal value, not a register.                                                                           |

## CPU Registers

There are 12 CPU registers, 8 of which are general purpose.

| Number | Alias | Purpose                                                                                                     |
| ------ | ----- | ----------------------------------------------------------------------------------------------------------- |
| `0x0`  | `ra`  | General purpose; used for whatever you desire.                                                              |
| `0x1`  | `rb`  | General purpose.                                                                                            |
| `0x2`  | `rc`  | General purpose.                                                                                            |
| `0x3`  | `rd`  | General purpose.                                                                                            |
| `0x4`  | `re`  | General purpose.                                                                                            |
| `0x5`  | `rf`  | General purpose.                                                                                            |
| `0x6`  | `rg`  | General purpose.                                                                                            |
| `0x7`  | `rh`  | General purpose.                                                                                            |
| `0x8`  | `ex`  | Extra: holds extra arithmetic output values, used for register expansions (`(value)` syntax).               |
| `0x9`  | `ac`  | Accumulator: holds the output of most arithmetic operations.                                                |
| `0xa`  | `sp`  | Stack pointer: holds the address of the top location in memory of the stack.                                |
| `0xb`  | `pc`  | Program counter: holds the address of the next instruction in memory to be executed.                        |
| `0xc`  | `bi`  | Boolean index: set to `0xffff` if the last cmp was equal, else `0xfffe`.                                    |

## The Simple Virtual Assembler

The assembler reads a rudimentary assembly language and outputs a binary format called "svb".
See the [`sva` directory](https://github.com/tteeoo/svc/tree/main/sva) for documentation on writing in the assembly language and using the assembler.

See the [`svd` directory](https://github.com/tteeoo/svc/tree/main/svd) for using the debugger and the [`asm` directory](https://github.com/tteeoo/svc/tree/main/asm) for some example programs.

## Memory

Sections:
* VGA text buffer: `0x00`-`0x7d0`
* Stack: `0x7d1`-`0x8ff`
* Program (varies in size): `0x900`-`0xX`
* Heap (everything else): `0xX+1`-`0xffff`

Before the CPU starts execution, a few things are done in memory:
* The value `0xffff` is pushed onto the stack. It will be pulled off with the "main" subroutine's `ret` instruction. When the program counter is set to `0xffff` the virtual machine will stop.
* The command-line arguments are loaded into the heap. (The heap is just all of the memory that doesn't have a specific purpose.)
* The word at `0xfffd` is set to the number of command-line arguments.
* The word at `0xfffe` is set to the size of the command-line arguments (number of characters + null terminators).
* The word at `0xffff` is set to the address of the start of the heap.

## To Do

* More example programs and documentation.
* Better tests (ones that actually exist).
* Keyboard input.
* Virtual drive with simple filesystem.

## License

The contents of this repository are made available under the Unlicense, a public domain equivalent license.

See the [`LICENSE` file](https://github.com/tteeoo/svc/tree/main/LICENSE) for its text.
