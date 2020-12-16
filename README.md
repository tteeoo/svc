# Simple Virtual Computer (svc)

A fictional 16-bit virtual machine written in Go.

It uses word-based memory, where each memory address maps to a word-length value.
This means that while there are only 65,536 memory addresses,
it technically has 128K of memory since each address points to a 16-bit value instead of a byte.

This repository contains the virtual machine and an assembler to compile programs for it.

## The Simple Virtual Assembler (sva)

In the directory `sva` you'll find the source code for the assembler.

It reads a rudimentary assembly language and outputs an equally simple binary format called svb.

Usage: `sva <input file> [-o output file]`

Documentation for the assembly language is to come.

## To Do

* A package for parsing and creating svb files.
* Instructions for sub-routines.
* Instructions for branching.
* A clock package.
* The ability to run svb files.
* Example programs and documentation.
* Some sort of debugging mode.
* Frame buffer/Graphics mode.

## License

svc is made available under the Unlicense, a public domain equivalent license.
