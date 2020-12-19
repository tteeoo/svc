# Simple Virtual Assembler

The assembler for the Simple Virtual Computer.

It reads a rudimentary assembly language and outputs an equally simple binary format called "svb".

## Usage

```
sva <input file> [-o output file]
```

## The Assembly Language

Comments can be denoted with `;`.

Each line in an input file does one of four things:

#### Sources another file
```
. <path to another file>
```
Works as if the contents of the other file were directly inserted into the current file at this line.

Examples:
```
. foo.asm
. dir/bar.asm
. ../baz.asm
```

#### Defines a constant
```
<name> = <value>
```
The value can be a double quoted string, hex value, or a positive/negative integer.
This will store some value at the next unique address(es) in memory after the offset.
This address (or the first address in the case of strings) can later be used in your program with the `[name]` syntax.
It does not matter where in the input file this constant is defined, it used anywhere in the program.

Examples:
```asm
foo = 0x41
bar = "Hello, world!"
baz = 42
qux = -1337
```

#### Defines an instruction to be executed
```
<name> <operands>...
```
Refer to the main README file to view a table of instuction names their operands.
Operands can be a hex value, positive/negative integer, a register alias (see main README as well), constant address (`[name]`), or a subroutine address (`{name}`).

Examples:
```asm
foo = "Z"
ldr 2 [foo] ; loads "Z" into register 2

cpl ac 0xff ; copies 0xff into the accumulator
cpl 0 257   ; copies 257 into register 0
sub 0       ; subtracts register 0's value from the accumulator
cpl 3 -2    ; copies -2 into register 3
cmp ac 3    ; compares the value of the accumulator and register 3
```

#### Defines a subroutine
```
<name>:
```
This will define a new subroutine with all of the instructions below it, until the next one is defined.
Subroutines can be used by jump instructions like `jmp`, `jme`, or `jne` with the `{name}` syntax.
Every program should have a `main` subroutine. This is a special subroutine that is compiled so that it is the entrypoint to your program (where the CPU starts reading).

Examples:
```asm
; prints "A" if 0xff equals 0xff

main:
  cpl ac 0xff ; copies 0xff into the accumulator
  cpl 0 0xff  ; copies 0xff into register 0
  cmp ac 0    ; compares the value of the accumulator and register 0
  jme {print}

print:
  cpl 3 0x0f41 ; copies 0x0f41 into register 3
  str 0 3      ; stores the value held in register 3 at the first address in memory
  vga          ; re-draws the text buffer, printing "A" with white text and black background
```

## The Simple Virtual Binary Format

To be written.
