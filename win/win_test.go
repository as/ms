package win

import (
	"image"
	"os/exec"
	"testing"
)

func ck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s", err)
	}
}
func TestOpen(t *testing.T) {
	cmd := exec.Command("winver.exe")
	err := cmd.Start()
	ck(t, err)
	defer cmd.Process.Kill()
	var win Window
	for i := 0; i < 1024; i++ {
		win, err = Open(cmd.Process.Pid)
		if err == nil {
			break
		}
	}
	ck(t, err)
	for y := 0; y < 1024; y += 64 {
		for x := 0; x < 1024; x += 64 {
			r := image.Rect(1, 1, x, y)
			err := win.Reshape(r)
			ck(t, err)
			r2, err := win.Bounds()
			ck(t, err)
			if r != r2 {
				t.Logf("want %s have %s", r, r2)
			}
		}
	}
}
