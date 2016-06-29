package peb

// Peb returns the location of the Process Environment Block
// it only works on amd64 Windows systems for 64-bit processes.
func Peb() uintptr
