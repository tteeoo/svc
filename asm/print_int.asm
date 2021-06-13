; Prints a signed integer, first converting it to a string.

. lib/io.asm
. lib/string.asm

main:

  ; Copy the uint into ac.
  cpl ac 2352
  sub (3567)

  ; We will store the string form of the uint in the heap, so we store
  ;   a pointer to the beginning of the heap in aa.
  ldr ra (0xffff)

  ; Convert the uint to a string.
  cal {itoa}

  ; Print the string.
  cpl rb 0
  cal {print}
  vga, ret
