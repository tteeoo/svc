; This file is intended to be sourced by another.
; It defines subroutines for dealing with strings.

; Converts an unsigned integer to a string.
; ac = The uint.
; aa = Starting address to write the string to.
utoa:
  ; Pseudocode:
  ; dd = 10
  ; cc = 0
  ; stack.push(aa)
  ; do {
  ;    ac /= dd
  ;    stack.push(ex)
  ;    cc++
  ; } while (ac != 0)
  ; while (cc != 0) {
  ;     bb = stack.pop()
  ;     ac = bb
  ;     ac += dd
  ;     memory[aa] = ac
  ;     aa++
  ;     cc--
  ; }
  ; memory[aa] = 0
  ; aa = stack.pop()

  ; Initialize registers:
  ; dd = Divisor.
  ; cc = Count.
  ; We push (and later pop) aa because its original value should be preserved.
  ; (It will likely be used right after this subroutine in printing the string.)
  cpl dd 10
  cpl cc 0
  psh aa

  ; Push remainder onto the stack until the result is zero.
  &loop_utoa
  div dd
  psh ex
  inc cc
  cml ac 0
  gtn &loop_utoa

  ; Pull off the stack and store until cc is 0.
  &after_utoa 
  cpl dd 48
  &loop2_utoa
  cml cc 0
  gte &after2_utoa
  pop bb
  cop ac bb
  add dd
  str aa ac
  dec cc
  inc aa
  gto &loop2_utoa
  &after2_utoa

  ; Add the null terminator.
  str aa (0)
  pop aa
  ret
