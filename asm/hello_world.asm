; Note: I still need to make the parser check for spaces
text = "Hello,World!"

print:
  ; The address of the character in ram should be in register 1
  ; The address to store the character in ram should be in register 2

  ; Load the character from ram
  ldr 0 1

  ; Apply the colors
  cop ac 0
  cpl 3 0x0f00
  orr 3

  ; Store the character in the text buffer
  str 2 ac

  ret


main:

  ; Initialize registers 1 and 2 with the start of the string
  ; and the number of the character, respectively
  cpl 1 [text]
  cpl 2 0

  ; Push the program counter on to the stack
  psh pc

  ; Print a character
  cal {print}

  ; Get ready for the next print
  inc 1
  inc 2

  ; Check if we have printed 12 characters
  cpl 3 12
  cmp 2 3

  ; Note: I should probably implement a more effecient way of looping
  pop 0
  dec 0
  dec 0
  gtn 0

  vga

  ret
