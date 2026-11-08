package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	BLACK = color.RGBA{0, 0, 0, 255}
	RED   = color.RGBA{220, 30, 30, 255}
)

type Numeric interface {
	~int | ~float64
}

func pos[T Numeric](x, y T) complex128 {
	return complex(float64(x), float64(y))
}

func line(img *image.RGBA, start, end complex128, col color.Color) {
	x, y := real(start), imag(start)
	for x < real(end) && y < imag(end) {
		img.Set(int(x), int(y), col)
		x++
		y++
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for i := range img.Pix {
		img.Pix[i] = 255 // white
	}
	line(img, pos(1, 1), pos(100, 200), RED)

	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}
