# Simple Virtual Computer

Currently under development; this is a work in progress.

A 16-bit virtual machine written in Go.

It uses word-based memory, where each memory address maps to a word-length value.
This means that while there are only 65,536 memory addresses,
it technically has 128K of memory since each address points to a 16-bit value instead of a byte.
The reasoning for this is simple: simplicity. Most modern computer systems use bytes because they can be more flexible and efficient.
Those are not the goals of this project.

It has a VGA text mode where it starts reading the contents of memory from a specified address (using an 80x25 resolution)
translates the encoded colors into ANSI escape codes,
and prints the output to your terminal.

This repository contains the virtual machine, an assembler to compile programs for it, and a debugger for those programs.

## Instruction Set

| Opcode | Name | Operands                              | Description                                                                                                                                              |
| ------ | ---- | --------                              | -----------                                                                                                                                              |
| `0x00` | nop |                                        | Does nothing                                                                                                                                             |
| `0x01` | cop | `reg to copy to` `reg to copy from`    | Copies the value from one register to another                                                                                                            |
| `0x02` | cpl | `reg` `value`                          | Copies a literal value to a register                                                                                                                     |
| `0x03` | str | `reg holding addr` `reg holding value` | Stores the value from a register into memory at the address held in another register                                                                     |
| `0x04` | ldr | `reg to load to` `reg holding addr`    | Loads the value from memory at the address held in one register into another register                                                                    |
| `0x05` | add | `reg`                                  | Adds the value held in a register to the accumulator                                                                                                     |
| `0x06` | sub | `reg`                                  | Subtracts the value held in a register from the accumulator                                                                                              |
| `0x07` | twc | `reg`                                  | Sets a register to the two's complement of the value it holds                                                                                            |
| `0x08` | inc | `reg`                                  | Increases the value held in a register by one                                                                                                            |
| `0x09` | dec | `reg`                                  | Decreases the value held in a register by one                                                                                                            |
| `0x0A` | mul | `reg`                                  | Multiplies the accumulator with the value held in a register                                                                                             |
| `0x0B` | div | `reg`                                  | Divides the accumulator by the value held in a register, storing the modulus in the "ex" register                                                        |
| `0x0C` | dcv | `reg`                                  | Divides the accumulator by the value held in a register, taking into account two's complement negative numbers, storing the modulus in the "ex" register |
| `0x0D` | xor | `reg`                                  | Performs the "xor" operation on the accumulator with the value held in a register                                                                        |
| `0x0E` | and | `reg`                                  | Performs the bitwise "and" operation on the accumulator with the value held in a register                                                                |
| `0x0F` | orr | `reg`                                  | Performs the bitwise "or" operation on the accumulator with the value held in a register                                                                 |
| `0x10` | not | `reg`                                  | Inverts the value held in a register                                                                                                                     |
| `0x11` | shr | `reg` `value to shift by`              | Shifts the value held in a register to the right by a specified value                                                                                    |
| `0x12` | shl | `reg` `value to shift by`              | Shifts the value held in a register to the left by a specified value                                                                                     |
| `0x13` | vga |                                        | Prints the VGA text buffer to the screen                                                                                                                 |
| `0x14` | psh | `reg`                                  | Decreases the stack pointer and sets the top value of the stack to the value held in a register                                                          |
| `0x15` | pop | `reg`                                  | Stores the top value of the stack in a register and increases the stack pointer                                                                          |
| `0x16` | ret |                                        | Pops the program counter off of the stack                                                                                                                |
| `0x17` | cal | `addr`                                 | Pushes the program counter onto the stack and sets the program counter to an address                                                                     |
| `0x18` | cmp | `reg` `reg`                            | If the values of two registers are the same, the boolean index is set to `0xffff`, else `0xfffe`                                                         |
| `0x19` | cle | `addr`                                 | Equivalent to "cal", but only executes if the boolean index is set to `0xffff`                                                                           |
| `0x1a` | cln | `addr`                                 | Equivalent to "cal", but only executes if the boolean index is set to `0xfffe`                                                                           |
| `0x1b` | gto | `addr`                                 | Sets the program counter to an address                                                                                                                   |
| `0x1c` | gte | `addr`                                 | Equivalent to "gto", but only executes if the boolean index is set to `0xffff`                                                                           |
| `0x1d` | gtn | `addr`                                 | Equivalent to "gto", but only executes if the boolean index is set to `0xfffe`                                                                           |
| `0x1f` | sth | `reg to load to` `reg holding addr`    | Like "sth", but offsets the address to make it store in the heap section of memory                                                                       |
| `0x1e` | ldh | `reg holding addr` `reg holding value` | Like "ldh", but offsets the address to make it load from the heap section of memory                                                                      |

## CPU Registers

| Number | Alias | Purpose                                                                                                     |
| ------ | ----- | -------                                                                                                     |
| `0`    | aa    | General purpose: used for whatever your program desires                                                     |
| `1`    | bb    | General purpose                                                                                             |
| `2`    | cc    | General purpose                                                                                             |
| `3`    | dd    | General purpose                                                                                             |
| `4`    | ex    | Extra: holds extra arithmetic output values, used for register expansions (`(value)` syntax)                |
| `5`    | ac    | Accumulator: holds the output of most arithmetic operations                                                 |
| `6`    | sp    | Stack pointer: holds the address of the top location in memory of the stack                                 |
| `7`    | pc    | Program counter: holds the address of the next instruction in memory to be executed                         |
| `8`    | bi    | Boolean index: set to `0xffff` if the last cmp was equal, else `0xfffe`                                     |

## The Simple Virtual Assembler

It reads a rudimentary assembly language and outputs a binary format called "svb".

See the [`sva` directory](https://github.com/tteeoo/svc/tree/main/sva) for more documentation on writing in the assembly language and using the assembler.
See the [`svd` directory](https://github.com/tteeoo/svc/tree/main/svd) for using the debugger.
See the [`asm` directory](https://github.com/tteeoo/svc/tree/main/asm) for some example programs.

## To Do

* Better tests (ones that actually exist).
* Keyboard input.
* More example programs and documentation.

## License

The contents of this repository are made available under the Unlicense, a public domain equivalent license.
