; prints "A" if 5 + 1 equals 6
foo = "A"

print:
  ldr ac ([foo]) ; loads "A" into the accumulator
  orr (0x0f00)   ; applies black background and white foreground colors to the accumulator
  str (0) ac     ; stores the value held in the accumulator at the first address in memory
  vga, ret       ; draws the text buffer, printing "A" with white text and black background
  
main:
  cpl ac 1    ; copies 1 into the accumulator
  add (5)     ; adds 5 to the accumulator
  cmp ac (6)  ; compares the value of the accumulator and 6
  cle {print} ; calls the print subroutine if the last cmp was equal
