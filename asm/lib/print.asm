; This file is intended to be sourced by another.
; It defines some routines for printing.

; Prints a character with a black background and white foreground.
print_char:
  ; The character to print should be in the accumulator.
  ; The address to store the character in ram should be in bb.

  ; Apply the colors (black background, white foreground).
  orr (0x0f00)

  ; Store the character in the text buffer.
  str bb ac

  ret

; Prints a null-terminated string using print_char.
print_str:
  ; The address of the start of the string should be in aa.
  ; The starting address to print to should be in bb.

  ; The following code is documented in the hello_world.asm example.
  ; It is duplicated there in order to use hello world as
  ; a more explanatory example, but featured here
  ; in order to be used in other programs.
  psh lc

  ldr ac aa

  cmp ac (0)

  cln {print_char}

  inc aa, bb

  pop, gtn cc

  ret
