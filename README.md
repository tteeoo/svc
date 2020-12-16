# Simple Virtual Computer (svc)

A fictional 16-bit virtual machine written in Go.

It uses word-based memory, where each memory address maps to a word-length value.
This means that while there are only 65,536 memory addresses,
it technically has 128K of memory since each address points to a 16-bit value instead of a byte.

It has a VGA text-mode where it starts reading the contents of memory from a specified address
(using a specified resolution, 80x25 by default),
translates the encoded colors into ANSI escape codes,
and prints the output to your terminal.

This repository contains the virtual machine and an assembler to compile programs for it.

## The Simple Virtual Assembler (sva)

In the directory `sva` you'll find the source code for the assembler.

It reads a rudimentary assembly language and outputs an equally simple binary format called "svb".

Usage: `sva <input file> [-o output file]`

Documentation for the assembly language will be made soon (see "To Do").

## Instruction Set

| Opcode  | Name  | Operands | Description  |
| ------- | ----- | -------- | ------------ |
| `0x00`  | nop |          | Does nothing |
| `0x01`  | cop | `register to copy to` `register to copy from` | Copies the value of one register to another |
| `0x02`  | cpl | `register to copy to` `value` | Copies a literal value into a register |
| `0x03`  | str | `register holding address` `register holding value` | Stores the value from a register into memory at the address held in another |
| `0x04`  | ldr | `register to load to` `register holding address` | Loads the value from memory at the address held in one register into another |
| `0x05`  | add | `register` | Adds the value of a register to the accumulator |
| `0x06`  | sub | `register` | Subtracts the value of a register from the accumulator |
| `0x07`  | twc | `register` | Sets a register to its two's complement |
| `0x08`  | inc | `register` | Increases the value held in a register by one |
| `0x09`  | dec | `register` | Decreases the value held in a register by one |
| `0x0A`  | mul | `register` | Multiplies the value held in a register with the accumulator |
| `0x0B`  | div | `register` | Divides the accumulator by the value held in a register, storing the modulus in the "ex" register |
| `0x0C`  | dcv| `register` | Divides the accumulator by the value held in a register, taking into account two's complement negative numbers, storing the modulus in the "ex" register |
| `0x0D`  | xor | `register` | Performs the "xor" operation on the accumulator and the value held in a register |
| `0x0E`  | and | `register` | Performs the bitwise "and" operation on the accumulator and the value held in a register |
| `0x0F`  | orr | `register` | Performs the bitwise "or" operation on the accumulator and the value held in a register |

## To Do

* A package for parsing and creating svb files.
* Instructions for sub-routines.
* Instructions for branching.
* A clock package.
* The ability to run svb files directly.
* Example programs and documentation.
* Some sort of debugging mode.
* Frame buffer/graphics mode.

## License

The contents of the repository are made available under the Unlicense, a public domain equivalent license.
