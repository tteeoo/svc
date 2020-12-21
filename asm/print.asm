; This file is intended to be sourced by another.
; It defines some routines for printing.

; Prints a character with a black background and white foreground.
print_char:
  ; The character to print should be in the accumulator.
  ; The address to store the character in ram should be in bb.

  ; Apply the colors (black background, white foreground).
  cpl dd 0x0f00
  orr dd

  ; Store the character in the text buffer.
  str bb ac

  ret
