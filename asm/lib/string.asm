; This file is intended to be sourced by another.
; It defines subroutines for dealing with strings.

; Converts an unsigned integer to a string.
; ac = The uint.
; ra = Starting address to write the string to.
utoa:
  ; Pseudocode:
  ; rd = 10
  ; rc = 0
  ; stack.push(ra)
  ; do {
  ;    ac /= rd
  ;    stack.push(ex)
  ;    rc++
  ; } while (ac != 0)
  ; while (rc != 0) {
  ;     rb = stack.pop()
  ;     ac = rb
  ;     ac += rd
  ;     memory[ra] = ac
  ;     ra++
  ;     rc--
  ; }
  ; memory[ra] = 0
  ; ra = stack.pop()

  ; Initialize registers:
  ; rd = Divisor.
  ; rc = Count.
  ; We push (and later pop) ra because its original value should be preserved.
  ; (It will likely be used right after this subroutine in printing the string.)
  cpl rd 10
  cpl rc 0
  psh ra

  ; Push remainder onto the stack until the result is zero.
  &loop_utoa
  div rd
  psh ex
  inc rc
  cml ac 0
  gtn &loop_utoa

  ; Pull off the stack and store until rc is 0.
  &after_utoa 
  cpl rd 48
  &loop2_utoa
  cml rc 0
  gte &after2_utoa
  pop rb
  cop ac rb
  add rd
  str ra ac
  dec rc
  inc ra
  gto &loop2_utoa
  &after2_utoa

  ; Add the null terminator.
  str ra (0)
  pop ra
  ret
