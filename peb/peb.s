//+build windows
//+build amd64

#include "textflag.h"

// func Peb() uintptr
TEXT ·Peb(SB),NOSPLIT,$0
	MOVQ	0x60(GS), AX
	MOVQ	AX, ret+0(FP)
	RET
