# Simple Virtual Assembler

The assembler for the Simple Virtual Computer.

It reads a rudimentary assembly language and outputs a file in a binary format called "svb".

How it works is a bit different than a typical x86 assembly language (partial because I don't know any "real" assembly languages very well).

## Usage

```
sva <input file> [-o <output file>] [-p]
```
`<output file>` will default to `./out.svb`.

With the `-p` option the assembler will write the preprocessed assembly to `<output file>.asm`.
Preprocessing includes stripping trailing whitespace and comments, sourcing files, and expanding instructions.
It can be useful for debugging.

To execute the assembled program, run:
```
svc <svb file>
```

## The Assembly Language

Comments are be denoted with `;`.

Each line of an input file does one of five things:

### Source another file
```
. <path to another file>
```
This works as if the contents of the other file were directly inserted into the current file at this line.
This will only go one file deep, i.e. a file that is sourced cannot source another file.

Examples:
```asm
. foo.asm
. dir/bar.asm
. ../baz.asm
```

### Define a constant
```
<name> = <value>
```
The value of the constant can be a double quoted string, hex (prefixed with `0x`), or a positive/negative integer.
This will store some value at a unique address in memory.
Negative integers will be stored in two's complement form, and strings will allocate each character consecutively, followed by a null word.
This address (or the address of the first character in the case of strings) can later be used in your program with the `[name]` syntax.

Examples:
```asm
foo = 0x41
bar = "Hello, world!"
baz = 42
qux = -1337
```

### Define an instruction to be executed
```
<name> <operands>...
```
Refer to the main `README.md` file to view a table of instruction names and their operands.
Operands can be a hex value, positive/negative integer, a register alias (also see main `README.md`), constant address (`[name]`), a subroutine address (`{name}`), or a label reference (`&name`).

Examples:
```asm
cpl ac 0xff ; Copies 0xff into the accumulator.
cpl aa 257  ; Copies 257 into register 0.
sub aa      ; Subtracts the value held in aa from the accumulator.
cpl dd -2   ; Copies -2 into register 3.
cmp ac dd   ; Compares the value of the accumulator and dd.
```

#### Instruction expansions
Expansions are syntactic sugar, allowing two instructions to be defined with one line.
There are three different types of them.

One type maps a single operation to two sets of operands. 
For example, this code:
```asm
inc aa, bb
```
expands to this:
```asm
inc aa
inc bb
```

Another variant maps two operations to one set of operands.
This code:
```asm
inc, add aa
```
will expand to:
```asm
inc aa
add aa
```

The last kind will take a value from inside some parenthesis, copy it into `ex`, and replace the value with "ex", so this code:
```asm
orr (0x0f00)
```
expands to:
```asm
cpl ex 0x0f00
orr ex
```

### Define a subroutine
```
<name>:
```
This will define a new subroutine with all of the instructions below it, until the next one is defined.
Subroutines' main purpose is to be used by call instructions (`cal`, `cle`, `cln`) with the `{name}` syntax.
Every program needs a "main" subroutine. This is a special subroutine that is compiled so that it is the entry point to your program (where the CPU starts executing).

Example:
```asm
; Prints "A" if 5 + 1 equals 6.
foo = "A"

print:
  ldr ac ([foo]) ; Loads "A" into the accumulator.
  orr (0x0f00)   ; Applies black background and white foreground colors to the accumulator.
  str (0) ac     ; Stores the value held in the accumulator at the first address in memory.
  vga, ret       ; Draws the text buffer, printing "A" with white text and black background.
  
main:
  cpl ac 1    ; Copies 1 into the accumulator.
  add (5)     ; Adds 5 to the accumulator.
  cmp ac (6)  ; Compares the value of the accumulator and 6.
  cle {print} ; Calls the print subroutine if the last cmp was equal.
  ret
```

There should be a `ret` instruction at the end of each subroutine, or else the CPU will continue executing whatever is stored after the subroutine in memory.
Subroutines cannot be addressed before they are defined.

### Define a label
```
&<name>
```
A label is like a subroutine except it can be addressed in code before it is defined, to make a "goto"-like structure.

Here is an example of a subroutine from `asm/lib/io.asm` that uses labels.
```asm
; Prints a string.
; aa = Address of the start of the string.
; bb = Starting address to print to.
print:

  ; The following code is like a while loop.
  ; Pseudo-code:
  ; while ((ac = memory[aa]) != 0) {
  ;     memory[bb] = ac | 0x0f00
  ;     aa++
  ;     bb++
  ; }

  ; Label the start of the subroutine for looping purposes.
  &loop_print_str

  ; If the value stored at the address in aa is 0x0, skip to the end.
  ldr ac aa
  cmp ac (0)
  gte &after_print_str

  ; Store the loaded character with the applied VGA color codes
  ;   at the address stored in bb, then increase aa and bb by 1.
  orr (0x0f00)
  str bb ac
  inc aa, bb

  ; Loop back up.
  gto &loop_print_str

  ; Label the end of the subroutine so it can be skipped to.
  &after_print_str
  ret
```
