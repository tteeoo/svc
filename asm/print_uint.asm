; Prints an unsigned integer, first converting it to a string.

. lib/io.asm
. lib/string.asm

main:

  ; Copy the uint into ac.
  cpl ac 0xfadb

  ; We will store the string form of the uint in the heap, so we store
  ;   a pointer to the beginning of the heap in aa.
  ldr ra (0xffff)

  ; Convert the uint to a string.
  cal {utoa}

  ; Print the string.
  cpl rb 0
  cal {print}
  vga, ret
