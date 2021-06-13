; This file is intended to be sourced by another.
; It defines subroutines for input and output.

; Prints a string.
; ra = Address of the start of the string.
; rb = Starting address to print to.
print:
  ; The following code is like a while loop.
  ; Pseudocode:
  ; while ((ac = memory[ra]) != 0) {
  ;     memory[rb] = ac | 0x0f00
  ;     ra++
  ;     rb++
  ; }

  ; Label the start of the subroutine for looping purposes.
  &loop_print_str

  ; If the value stored at the address in aa is 0x0, skip to the end.
  ldr ac ra
  cml ac 0
  gte &after_print_str

  ; Store the loaded character with the applied VGA color codes
  ;   at the address stored in bb, then increase aa and bb by 1.
  orr (0x0f00)
  str rb ac
  inc ra, rb

  ; Loop back up.
  gto &loop_print_str

  ; Label the end of the subroutine so it can be skipped to.
  &after_print_str
  ret
