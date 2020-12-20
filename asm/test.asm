; comment
bar = "AAA"
foo = 0x88
another:
	nop
	vga
	cpl ex 0xff
main:
	cpl ac -2
