; Define the text we want to print.
; Each character will be consecutively stored in memory,
; with the first address available via "[text]".
; The address after the end of the string will be set to 0
; to denote the end of the string.
text = "Hello, World!"

; Source the print.asm file.
; It contains the print_char subroutine.
. lib/print2.asm

main:

  ; Initialize registers aa and bb.
  ; aa will hold the address of the character to print in memory.
  ; bb will hold the number of characters that have been printed.
  cpl aa [text]
  cpl bb 0

  ; Push lc on to the stack.
  ; The value of lc is equivalent to the address
  ; of the current executing instruction.
  ; This will be pushed onto the stack and used to loop back to.
  cal {print_str}

  ; Draw the text buffer.
  vga

