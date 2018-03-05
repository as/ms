package win

import (
	"image"
)

// Rect returns the windows bounds on screen
func Rect() (image.Rectangle, bool) { return rect() }

// ClientRect returns the sub-rectangle of the client area (the intersection of
// the window border and the window). The bounds are relative to the window bounds,
// so this rectangle remains constant if the window is moved without resize
func ClientRect() (image.Rectangle, bool) {
	r, _ := clientRect()
	return r, false
}

// ClientAbs returns the absolute client area.
func ClientAbs() image.Rectangle {
	wr, _ := Rect()
	cr, _ := ClientRect()
	return wr.Add(Border(wr, cr))
}

func Border(wr, cr image.Rectangle) image.Point {
	return image.Point{
		X: wr.Max.X - cr.Max.X - wr.Min.X,
		Y: wr.Max.Y - cr.Max.Y - wr.Min.Y,
	}
}
func FromPID(p int) (wids []uintptr) { return fromPID(p) }
func Move(wid int, to image.Rectangle, paint bool) (err error) {
	return move(wid, to, paint)
}
