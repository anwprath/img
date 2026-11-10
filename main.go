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

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	ax, ay := 7, 3
	bx, by := 12, 37
	cx, cy := 62, 53

	img.Set(ax, ay, white)
	img.Set(bx, by, white)
	img.Set(cx, cy, white)


	flipVertical(img) // inverted coordinate system of Go who tf knows why
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}
