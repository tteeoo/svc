# Simple Virtual Assembler

The assembler for the Simple Virtual Computer.

It reads a rudimentary assembly language and outputs an equally simple binary format called "svb".

## Usage

```
sva <input file> [-o <output file>] [-p]
```
`<output file>` will default to `./out.svb`.

The `-p` option will write the pre-processed assembly to `<output file>.asm`.
Pre-processing includes stripping trailing whitespace and comments, sourcing files, and expanding instructions.
It can be useful for debugging.

To execute the assembled program, run:
```
svc <svb file>
```

## The Assembly Language

Comments can be denoted with `;`.

Each line in an input file does one of four things:

#### Sources another file
```
. <path to another file>
```
Works as if the contents of the other file were directly inserted into the current file at this line.
This will only go one file deep. A file that is sourced cannot source another file.

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
This will store some value at a unique address in memory.
Negative integers will become their two's complement, and strings will allocate each character consecutively, followed by a null terminator.
This address (or the first address in the case of strings) can later be used in your program with the `[name]` syntax.

Examples:
```[asm](asm)
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
cpl ac 0xff ; copies 0xff into the accumulator
cpl aa 257  ; copies 257 into register 0
sub aa      ; subtracts the value held in aa from the accumulator
cpl dd -2   ; copies -2 into register 3
cmp ac dd   ; compares the value of the accumulator and dd
```

##### Instruction expansions
An instruction expansion is syntactic sugar, allowing two instructions to be defined with one line.
They can take two forms: either mapping one operation to two sets of operands, or mapping two operations to one set of operands.

For example, this code:
```
inc aa, bb
```
Will expand to this:
```
inc aa
inc bb
```

And this code:
```
inc, add aa
```
Will expand to:
```
inc aa
add aa
```

#### Defines a subroutine
```
<name>:
```
This will define a new subroutine with all of the instructions below it, until the next one is defined.
Subroutines can be used by call instructions (`cal`, `cle`, `cln`) with the `{name}` syntax.
Every program needs a `main` subroutine. This is a special subroutine that is compiled so that it is the entrypoint to your program (where the CPU starts reading).

Examples:
```asm
; prints "A" if 0xff equals 0xff
foo = "A"

print:
  ldr ac [foo]  ; loads "A" into the accumulator
  orr 0x0f00    ; applies black background and white foreground colors to the accumulator
  cpl bb 0      ; copies 0 into bb
  str bb ac     ; stores the value held in the accumulator at the first address in memory
  vga, ret      ; draws the text buffer, printing "A" with white text and black background
  
main:
  cpl cc 0xff, aa 0xff ; copies 0xff into the accumulator and aa
  cmp cc aa            ; compares the value of cc and aa
  cle {print}          ; calls the print subroutine if 0xff equals 0xff
```

## The Simple Virtual Binary Format

To be written.
