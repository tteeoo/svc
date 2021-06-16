package cpu

// Pipelining will work as follows:
// Four processes for an instruction:
//   fetch+decode, read regs, execute, write mem, write regs.
// One type of struct per process, implementing an interface which has two methods:
//   Process(), and Promote() returns the next struct or nil.
// Use this algorithm on a four length array of instruction interfaces:
//   1. Add one interface to the array where there is nil
//   2. .Process() all instructions
//   3. .Promote() all instructions, setting index to the return value
// TODO: plan flushing, interrupts, stopping execution.
