main:
  cpl 0 0x0f41 ; copies 0x0f41 into register 0
  cpl 1 0      ; copies 0 into register 1
  str 1 0      ; stores the value held in register 0 at the first address in memory
  vga          ; draws the text buffer, printing "A" with white text and black background
