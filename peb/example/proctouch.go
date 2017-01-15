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
	"bytes"
	"os"
)

func wide(s []byte) []uint16{
	return utf16.Encode([]rune(string(s)))
}
func toslice(p *uint16) []uint16{
	return (*(*[65536]uint16) (unsafe.Pointer(p)))[:]
}

func dupargs(a []string) []string{
	b := make([]string, len(a))
	for i, v := range a{
		b[i] = v
	}
	return b
}

func bcopy(p peb.BStr, s string){
	x := toslice(p.P)
	copy(x, wide([]byte(s)))			
}

func main(){
	if len(os.Args) < 3{
		fmt.Println("usage: proctouch imagepath commandline newcwd")
	}
	Peb := peb.Peb
	args := dupargs(os.Args)
	path := args[1] + "\x00"
	cli  := args[2] + "\x00"
	cwd  := args[3] + "\x00"
	
	printPEB()

	// Change some things just because
	bcopy(Peb.Params.CWD, cwd)	
	bcopy(Peb.Params.CommandLine, cli)	
	bcopy(Peb.Params.ImagePathName, path)
	bcopy(Peb.Loader.Order.ByMemory.Next.FullDLL, "eggs")
	
	Peb.Session = 0
	
	os.Chdir("C:\\")
	envBlank()
	printPEB()
	fmt.Println("done; open process explorer and check pid", os.Getpid())
	nullify()
	time.Sleep(100*time.Second)
}

func printPEB(){
	fmt.Printf("%s", peb.Peb.String())
}

func nullify(){
	peb.Peb.Params = nil
	peb.Peb.Loader = nil
}

func envBlank(){
	v := "COMPUTERNAME"
	fmt.Println(v, os.Getenv(v))
	evp := (*(*[65535]uint16)(unsafe.Pointer(peb.Peb.Params.Env)))[:65535]
	copy(evp, wide(bytes.Repeat([]byte("M\x00I\x00N\x00K\x00"), 100)))
	fmt.Println(v, os.Getenv(v))
}

func stderrHijack(filename string) {
	file, _ := os.Create(filename)
	peb.Peb.Params.Stderr = file.Fd()
}