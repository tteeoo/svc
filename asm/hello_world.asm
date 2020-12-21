; Define the text we want to print.
; Each character will be consecutively stored in memory,
; with the first address available via "[text]".
; The address after the end of the string will be set to 0
; to denote the end of the string.
text = "Hello, World!"

; Source the print.asm file.
. print.asm

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
  psh lc

  ; Load the value at the address in aa into the accumulator.
  ldr ac aa

  ; Check if the accumulator is equal to 0, the "null terminator",
  ; signifying the end of a string.
  cpl dd 0
  cmp ac dd

  ; Print a character if it is not the null terminator.
  ; This will call the print_char subroutine defined in the print.asm file.
  cln {print_char}

  ; Get ready for the next print.
  ; Since the text we want to print is consecutively stored in memory
  ; we can just increment the address to get the next one.
  ; We also increment the number of characters printed.
  inc aa, bb

  ; Loop back to the psh instruction if the loaded character
  ; was not the null terminator.
  pop, gtn cc

  ; Draw the text buffer.
  vga

