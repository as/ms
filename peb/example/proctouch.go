//+build windows
//+build amd64
package main
// Given an imagepath and commandline, alter the imagepath and command line string
// of the current process

import(
	"github.com/as/ms/peb"
	"fmt"
	"unsafe"
	"unicode/utf16"
	"time"
	"os"
)

func wide(s []byte) []uint16{
	return utf16.Encode([]rune(string(s)))
}
func toslice(p *uint16) []uint16{
	return (*(*[65536]uint16) (unsafe.Pointer(p)))[:]
}

func main(){
	if len(os.Args) < 2{
		fmt.Println("usage: proctouch imagepath commandline")
	}
	path := make([]byte, len(os.Args[1])+1)
	cli  := make([]byte, len(os.Args[2])+1)
	copy(path, os.Args[1])
	copy(cli, os.Args[2])
	
	fmt.Printf("%s", peb.Peb.String())
	
	x := toslice(peb.Peb.Params.CommandLine.P)
	copy(x, wide(cli))
	x = toslice(peb.Peb.Params.ImagePathName.P)
	copy(x, wide(path))
	
	fmt.Printf("%s", peb.Peb.String())
	time.Sleep(2*time.Second)

	fmt.Println("done; open process explorer and check pid", os.Getpid())
	time.Sleep(100*time.Second)
}
