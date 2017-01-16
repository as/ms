package main

import (
	"flag"
	"github.com/as/ms/win"
	"image"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var random = flag.Bool("r", false, "reshape window to randomly generated rectangle")

func init() {
	rand.Seed(time.Now().Unix())
}

func genrandom() image.Rectangle {
	x := rand.Intn(1920)
	y := rand.Intn(1080)
	dx := rand.Intn(1920 - x)
	dy := rand.Intn(1080 - y)
	return image.Rect(x, y, dx, dy)
}

func atoi(a string) (i int) {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}
func main() {
	flag.Parse()
	A := flag.Args()
	wid := win.FromPID(atoi(A[0]))
	var to image.Rectangle
	if !*random {
		to = image.Rect(atoi(A[1]), atoi(A[2]), atoi(A[3]), atoi(A[4]))
	} else {
		to = genrandom()
	}
	if len(wid) == 0 {
		log.Fatalf("pid %d has no windows\n")
	}
	for {
		for _, w := range wid {
			win.Move(int(w), to, true)
			if !*random {
				break
			}
			to = genrandom()
		}
	}
}
