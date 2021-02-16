; Define the text we want to print.
; Each character will be consecutively stored in memory,
;   with the first character addressable by "[text]".
; The address after the end of the string will be set to 0x0 to denote the end
;   of the string.
text = "Hello, World!"

; Source the io.asm file.
; It contains the print subroutine.
. lib/io.asm

; The "main" subroutine is the entry point of the program.
main:

  ; Initialize registers aa and bb.
  ; aa holds the address of the character in memory to print.
  ; bb holds the address of where in memory to print to.
  cpl aa [text]
  cpl bb 0

  ; Call the print subroutine.
  cal {print}

  ; Draw the text buffer.
  vga

  ; There is a "ret" instruction automatically added after the main subroutine.

