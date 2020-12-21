; Note: I still need to make the parser check for spaces.
text = "Hello,World!"

print:
  ; The address of the character in ram should be in aa.
  ; The address to store the character in ram should be in bb.

  ; Load the character from ram.
  ldr ac aa

  ; Apply the colors (black background, white foreground).
  cpl dd 0x0f00
  orr dd

  ; Store the character in the text buffer.
  str bb ac

  ; Return back to after the cal instruction.
  ret

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

  ; Print a character.
  cal {print}

  ; Get ready for the next print.
  ; Since the text we want to print is consecutively stored in memory
  ; we can just increment the address to get the next one.
  ; We also increment the number of characters printed.
  inc aa
  inc bb

  ; Check if we have printed 12 characters.
  cpl dd 12
  cmp bb dd

  ; Loop back to the psh instruction if we have not printed 12 characters.
  pop cc
  gtn cc

  ; Draw the text-buffer.
  vga

