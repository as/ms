package peb

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Peb points to the PEB in memory. This pointer is a fixed offset into
// userspace virtual memory. This offset is decided some time upon OS init.
// Address-wise, all processes see the offset for Peb, but the address itself
// points to a different physical address location depending on which process
// is accessing the offset.
//
// To alter another proces's PEB, access uintptr(Peb) in the target
// process's address space
var Peb *PEB

// Peb returns the location of the Process Environment Block
// it only works on amd64 Windows systems for 64-bit processes.
func getpeb() uintptr

func init() {
	Peb = (*PEB)(unsafe.Pointer(getpeb()))
}

// PEB is the memory layout of the process environment block. It describes the
// state of the process environment and, through theory and practice, has proven
// to be mutable by the running process as well as any other process running
// under a similar security context.
//
// So far, my observation is that any process running as the same user
// may alter this structure without requesting special privledges
//
// Most of the structures in the PEB are mutable
type PEB struct {
	_        [2]byte
	Debugged byte
	_        [1]byte
	_        [2]byte
	_        [18]byte // amd64 only?
	Loader   *Loader
	Params   *Params
	_        [104]byte
	_        [52]uintptr
	_        uintptr
	_        [128]byte
	_        uintptr
	Session  uint
}

// Loader contains the state of Window's image loader. It mainly serves as
// a pointer to the head of a linked list of table entries.
type Loader struct {
	_       [8]byte
	_       [1]uintptr // if no work try 3
	Next, _ *TableEntry
}

// TableEntry contains information about loaded modules. This information
// is owned by the running process and remains mutable throughout the process's
// lifetime.
type TableEntry struct {
	_             [2]uintptr
	Next, Prev    *TableEntry
	_             [2]uintptr
	DLLBase       uintptr
	EntryPoint    uintptr
	_             [1]uintptr
	FullDLLName   LPWStr
	_             [8]byte
	_             [3]uintptr
	CheckSum      uint
	TimeDateStamp uint
}

// Params contains process startup state.
type Params struct {
	_ [16]byte
	_ [11]uintptr

	// This is the path to the executable. Altering this
	// confuses naive anti-virus software and even Mark Russinovich's
	// "process monitor".
	//
	// For interesting behavior, change it to a remote share and look
	// at the resulting network traffic on ports 139 or 445.
	//
	ImagePathName LPWStr

	// The 'Unknown' contains the ImagePath's length
	// my guess is it also contains information about
	// the command line arguments.
	//
	// Mutant COM Bstr?
	//
	Unknown uintptr

	// Mutable argument vector
	CommandLine LPWStr
}

// LPWstr is a pointer to a wide string. Internally, P either points
// to the starting index of a []uint16, or is nil.
type LPWStr struct {
	P *uint16
}

func (t PEB) String() string {
	s := fmt.Sprintf("===PEB===\n")
	s += fmt.Sprintln(t.Loader)
	s += fmt.Sprintln(t.Params)
	s += fmt.Sprintln("Session", t.Session)
	return s
}
func (t Loader) String() string {
	s := fmt.Sprintf("\n===Loader===\n")
	s += fmt.Sprint("Entry:")
	if t.Next != nil {
		s += fmt.Sprintln(t.Next)
	} else {
		s += fmt.Sprintln("<nil>")
	}
	return s
}
func (t TableEntry) String() string {
	s := fmt.Sprintf("	%#v", t)
	if t.Next != nil {
		s += fmt.Sprintf("\n%s", t.Next)
	}
	return s
}
func (p Params) String() string {
	s := fmt.Sprintf("===Params===\n")
	s += fmt.Sprintf("Unknown: %0x\n", p.Unknown)
	s += fmt.Sprintf("ImagePathName: %q\n", p.ImagePathName)
	s += fmt.Sprintf("CommandLine: %q\n", p.CommandLine)
	return s
}

func (l LPWStr) String() string {
	if l.P != nil {
		return syscall.UTF16ToString((*(*[65535]uint16)(unsafe.Pointer(l.P)))[:])
	}
	return "<empty>"
}
