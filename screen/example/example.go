package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"image"
	"github.com/as/screen"
)

func main() {
	var err error
	scr, win := 0, 0
	l := len(os.Args)
	if l > 1 {
		scr, err = strconv.Atoi(os.Args[1])
		no(err)
	}
	if l > 2 {
		win, err = strconv.Atoi(os.Args[2])
		no(err)
	}

	fmt.Println(screen.Capture(scr, win, image.Rect(0, 0, 1024, 768)))
}

func no(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
