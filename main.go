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
	int | float64 | int32
}

func abs[T Numeric](x T) T {
	if x > 0 {
		return x
	}
	return -x
}

func line(img *image.RGBA, a, b int32, c, d int32, col color.Color) {
	dx, dy := abs(a-c), -abs(b-d)
	sx, sy := 1, 1

	if a > c {
		sx = -1
	}
	if b > d {
		sy = -1
	}
	err := dx + dy
	x0, y0 := int(a), int(b)
	x1, y1 := int(c), int(d)
	for {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 10000, 10000))
	for i := range img.Pix {
		img.Pix[i] = 255 // white
	}
	line(img, 0, 0, 10000, 2000, RED)
	line(img, 0, 0, 10000, 2001, RED)
	line(img, 0, 0, 10000, 2002, RED)
	line(img, 0, 0, 10000, 2003, RED)
	line(img, 0, 0, 10000, 2004, RED)
	line(img, 0, 0, 10000, 2005, RED)
	line(img, 0, 0, 10000, 2006, RED)
	line(img, 0, 0, 10000, 2007, RED)
	line(img, 0, 0, 10000, 2008, RED)
	line(img, 0, 0, 10000, 2009, RED)
	line(img, 0, 0, 10000, 2010, RED)

	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}
