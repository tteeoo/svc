# Simple Virtual Computer

Currently under development; this is a work in progress.

A 16-bit virtual machine written in Go.

It uses word-based memory, where each memory address maps to a word-length value.
This means that while there are only 65,536 memory addresses,
it technically has 128K of memory since each address points to a 16-bit value instead of a byte.
The reasoning for this is simple: simplicity. Most modern computer systems use bytes because they can be more flexible and effecient.
Those are not the goals of this project.

It has a VGA text mode where it starts reading the contents of memory from a specified address
(using a specified resolution, 80x25 by default),
translates the encoded colors into ANSI escape codes,
and prints the output to your terminal.

This repository contains the virtual machine and an assembler to compile programs for it.

## Instruction Set

| Opcode | Name | Operands                                            | Description                                                                                                                                              |
| ------ | ---- | --------                                            | -----------                                                                                                                                              |
| `0x00` | nop  |                                                     | Does nothing                                                                                                                                             |
| `0x01` | cop  | `register to copy to` `register to copy from`       | Copies the value from one register to another                                                                                                            |
| `0x02` | cpl  | `register` `value`                                  | Copies a literal value to a register                                                                                                                     |
| `0x03` | str  | `register holding address` `register holding value` | Stores the value from a register into memory at the address held in another register                                                                     |
| `0x04` | ldr  | `register to load to` `register holding address`    | Loads the value from memory at the address held in one register into another register                                                                    |
| `0x05` | add  | `register`                                          | Adds the value held in a register to the accumulator                                                                                                     |
| `0x06` | sub  | `register`                                          | Subtracts the value held in a register from the accumulator                                                                                              |
| `0x07` | twc  | `register`                                          | Sets a register to the two's complement of the value it holds                                                                                            |
| `0x08` | inc  | `register`                                          | Increases the value held in a register by one                                                                                                            |
| `0x09` | dec  | `register`                                          | Decreases the value held in a register by one                                                                                                            |
| `0x0A` | mul  | `register`                                          | Multiplies the accumulator with the value held in a register                                                                                             |
| `0x0B` | div  | `register`                                          | Divides the accumulator by the value held in a register, storing the modulus in the "ex" register                                                        |
| `0x0C` | dcv  | `register`                                          | Divides the accumulator by the value held in a register, taking into account two's complement negative numbers, storing the modulus in the "ex" register |
| `0x0D` | xor  | `register`                                          | Performs the "xor" operation on the accumulator with the value held in a register                                                                        |
| `0x0E` | and  | `register`                                          | Performs the bitwise "and" operation on the accumulator with the value held in a register                                                                |
| `0x0F` | orr  | `register`                                          | Performs the bitwise "or" operation on the accumulator with the value held in a register                                                                 |
| `0x10` | not  | `register`                                          | Inverts the value held in a register                                                                                                                     |
| `0x11` | shr  | `register` `value to shift by`                      | Shifts the value held in a register to the right by a specified value                                                                                    |
| `0x12` | shl  | `register` `value to shift by`                      | Shifts the value held in a register to the left by a specified value                                                                                     |
| `0x13` | vga  |                                                     | Prints the VGA text buffer to the screen                                                                                                                 |
| `0x14` | psh  | `register`                                          | Decreases the stack pointer and sets the top value of the stack to the value held in a register                                                          |
| `0x15` | pop  | `register`                                          | Stores the top value of the stack in a register and increases the stack pointer                                                                          |
| `0x16` | ret  |                                                     | Pops the program counter off of the stack                                                                                                                |
| `0x17` | cal  | `address`                                           | Pushes the program counter onto the stack and sets the program counter to an address                                                                     |
| `0x18` | cmp  | `register` `register`                               | If the values of two registers are the same, the extra register is set to `0xffff`, else `0xfffe`                                                        |
| `0x19` | cle  | `address`                                           | Equivalent to "cal", but only executes if the extra register is set to `0xffff`                                                                          |
| `0x1a` | cln  | `address`                                           | Equivalent to "cal", but only executes if the extra register is set to `0xfffe`                                                                          |
| `0x1b` | gto  | `register`                                          | Sets the program counter to the value held in a register                                                                                                 |
| `0x1c` | gte  | `register`                                          | Equivalent to "gto", but only executes if the extra register is set to `0xffff`                                                                          |
| `0x1d` | gtn  | `register`                                          | Equivalent to "gto", but only executes if the extra register is set to `0xfffe`                                                                          |

## CPU Registers

| Number | Alias | Purpose                                                                                  |
| ------ | ----- | -------                                                                                  |
| `0`    | aa    | General purpose: used for whatever your program desires                                  |
| `1`    | bb    | General purpose                                                                          |
| `2`    | cc    | General purpose                                                                          |
| `3`    | dd    | General purpose                                                                          |
| `4`    | ex    | Extra: holds the output of miscellaneous instructions, or extra arithmetic output values |
| `5`    | ac    | Accumulator: holds the output of most arithmetic operations                              |
| `6`    | sp    | Stack pointer: holds the address of the top location in memory of the stack              |
| `7`    | pc    | Program counter: holds the addess of the next instruction in memory to be executed       |
| `8`    | lc    | Last counter: holds the last value of the program counter, useful for loops              |

## The Simple Virtual Assembler (sva)

In the directory `sva` you'll find the source code for the assembler.

It reads a rudimentary assembly language and outputs an equally simple binary format called "svb".

Usage:
```
sva <input file> [-o output file]
```

Then, to run the assembled program, run:
```
svc <svb file>
```

See the ["sva" directory](https://github.com/tteeoo/svc/tree/main/sva) for more documentation, and the ["asm" directory](https://github.com/tteeoo/svc/tree/main/asm) for some example programs.

## To Do

* Better tests (ones that actually exist).
* Keyboard input.
* Some sort of debugging mode.
* More example programs and documentation.

## License

The contents of this repository are made available under the Unlicense, a public domain equivalent license.
