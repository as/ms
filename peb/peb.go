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
	SharedAddr byte
	ExecOptions byte
	Debugged byte
	_        byte
	//_        [2]byte
	//_        [18]byte // amd64 only?
	BaseAddr uintptr
	Mutex     uintptr
	
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
	Size    uint32
	Done    uint32
	_       [1]uintptr // if no work try 3
	Order
}

type Order struct{
	ByLoad	Link
	ByMemory	Link
	ByInit	Link
}

type Link struct{
	Next, Prev *Module
}

// Module contains information about loaded modules. This information
// is owned by the running process and remains mutable throughout the process's
// lifetime.
type Module struct {
	Order
	DLLBase       uintptr
	EntryPoint    uintptr
	ImageSize     uint32
	FullDLL   BStr
	BaseDLL   BStr
	Flags     uint
	NLoaded	  uint16
	TLSIndex  uint16
	Section   uintptr
	CheckSum      uint
	TimeDateStamp uint
	Ctx uintptr
	_ uintptr
	_ [3]uintptr
}

// Params contains process startup state.
type Params struct {
	AllocSize uint32
	Size uint32
	Flags uint32
	DebugFlags uint32
	//_ [16]byte

	Console uintptr
	PGroup uint32
	Stdin, Stdout, Stderr uintptr // careful, these arent 0, 1, and 2
	CWD BStr
	CWDHandle uintptr
	DLLPath BStr
	
	//_ [11]uintptr

	// This is the path to the executable. Altering this
	// confuses naive anti-virus software and even Mark Russinovich's
	// "process monitor".
	//
	// For interesting behavior, change it to a remote share and look
	// at the resulting network traffic on ports 139 or 445.
	//
	ImagePathName BStr
	
	// Mutable argument vector
	CommandLine BStr
	
	// Supposedly points to the environment, but probabilistic crashes
	// occur when dereferencing the pointer inside...
	Env uintptr	

	
}

type BStr struct{
	Len uint16
	Cap uint16
	LPWStr
}
// LPWstr is a pointer to a wide string. Internally, P either points
// to the starting index of a []uint16, or is nil.
type LPWStr struct {
	P *uint16
}

func (t PEB) String() string {
	s := fmt.Sprintf("===PEB===\n")
	s  = fmt.Sprintf("=Loader=\n\n")
	s += fmt.Sprintln(t.Loader)
	
	s  = fmt.Sprintf("\n=Params=\n\n")
	s += fmt.Sprintln(t.Params)
	
	s += fmt.Sprintln("Session", t.Session)
	return s
}
func (t Loader) String() string {
	return fmt.Sprintln(&(t.Order))
}
func (l *Order) String() string{
	s := ""
	mod := l.ByMemory.Next
	var pmod *Module
	for i := 0; i < 2; i++{
		fmt.Printf("sp=%p next=%p\n",  pmod, mod)
		s += fmt.Sprintln(mod)
		pmod = mod
		mod = mod.Order.ByMemory.Next
	}
	return s
}

func (m *Module) String() string{
	fmt.Printf("%q\n", m.FullDLL)
	if m == nil{
		return fmt.Sprint("<nil>")
	}
	return fmt.Sprintf("%#v", m.FullDLL)
}

func (p Params) String() string {
	s := fmt.Sprintf("===Params===\n")
	//s += fmt.Sprintf("Unknown: %0x\n", p.Unknown)
	s += fmt.Sprintf("CWD: %q\n", p.CWD)
	s += fmt.Sprintf("DLLPath: %q\n", p.DLLPath)
	s += fmt.Sprintf("ImagePathName: %q\n", p.ImagePathName)
	s += fmt.Sprintf("CommandLine: %q\n", p.CommandLine)
	// s += fmt.Sprintf("Env: 0x%0x\n", uintptr(unsafe.Pointer(p.Env.P)))
	return s
}

func (b BStr) String() string{
	return fmt.Sprintf("(len=%d cap=%d) %s", b.Len, b.Cap, b.LPWStr)
}

func (l LPWStr) String() string {
	if l.P != nil {
		return syscall.UTF16ToString((*(*[65535]uint16)(unsafe.Pointer(l.P)))[:])
	}
	return "<empty>"
}
