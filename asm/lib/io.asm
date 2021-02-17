; This file is intended to be sourced by another.
; It defines subroutines for input and output.

; Prints a string.
; aa = Address of the start of the string.
; bb = Starting address to print to.
print:
  ; The following code is like a while loop.
  ; Pseudocode:
  ; while ((ac = memory[aa]) != 0) {
  ;     memory[bb] = ac | 0x0f00
  ;     aa++
  ;     bb++
  ; }

  ; Label the start of the subroutine for looping purposes.
  &loop_print_str

  ; If the value stored at the address in aa is 0x0, skip to the end.
  ldr ac aa
  cml ac 0
  gte &after_print_str

  ; Store the loaded character with the applied VGA color codes
  ;   at the address stored in bb, then increase aa and bb by 1.
  orr (0x0f00)
  str bb ac
  inc aa, bb

  ; Loop back up.
  gto &loop_print_str

  ; Label the end of the subroutine so it can be skipped to.
  &after_print_str
  ret
