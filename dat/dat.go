// Package dat contains information about CPU opcodes and registers.
package dat

func init() {
	// Create the CodeToName map.
	for k, v := range OpNameToCode {
		OpCodeToName[v] = k
	}
}
