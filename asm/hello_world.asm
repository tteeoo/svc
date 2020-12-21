text = "Hello,World!"

main:

  cpl 1 [text]
  cpl 2 0

  ldr 0 1
  cop ac 0
  cpl 3 0x0f00
  orr 3
  str 2 ac

  inc 1
  inc 2

  ldr 0 1
  cop ac 0
  cpl 3 0x0f00
  orr 3
  str 2 ac

  inc 1
  inc 2

  vga
