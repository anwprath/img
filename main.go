package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	white  = color.RGBA{255, 255, 255, 255}
	green  = color.RGBA{0, 255, 0, 255}
	red    = color.RGBA{255, 0, 0, 255}
	blue   = color.RGBA{64, 128, 255, 255}
	yellow = color.RGBA{255, 200, 0, 255}
)

type Numeric interface {
	int | float64 | int32
}

func abs[T Numeric](x T) T {
	if x == 0 {
		return 1
	}
	if x > 0 {
		return x
	}
	return -x
}

// x = x0 + m*(x1-x0)
// y = y0 + m*(y1-y0)
func line(img *image.RGBA, c color.RGBA, x0, y0, x1, y1 int) {
	for x := x0; x <= x1; x += 1 {
		f := float64(x-x0) / float64(x1-x0)
		y := int(float64(y0) + f*(float64(y1)-float64(y0)))
		img.Set(x, y, c)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	ax, ay := 7, 3
	bx, by := 12, 37
	cx, cy := 62, 53

	line(img, blue, ax, ay, bx, by)
	line(img, green, cx, cy, bx, by)
	line(img, yellow, cx, cy, ax, ay)
	line(img, red, ax, ay, cx, cy)

	img.Set(ax, ay, white)
	img.Set(bx, by, white)
	img.Set(cx, cy, white)

	flipVertical(img) // inverted coordinate system of Go who tf knows why
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}
