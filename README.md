# Simple Virtual Computer (svc)

A 16-bit virtual computer written in Go.

It uses word-based memory, each memory address maps to a word-length value (instead of a byte).
This means that while it has 65,536 memory addresses, it has 128K of memory since each of those addresses has a 16 bit value instead of 8.

This repository contains the virtual machine and an assembler.

## License

svc is made available under the Unlicense, a public domain equivalent license.
