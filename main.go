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

func circle(img *image.RGBA, cx, cy, radius int, col color.Color) {
	x, y := radius, 0
	// err := 1 - radius

	for x >= y {
		octants := [][2]int{
			{x, y}, {x, -y}, {-x, y}, {-x, -y},
			{y, x}, {y, -x}, {-y, x}, {-y, -x},
		}
		for _, point := range octants {
			img.Set(cx+point[0], cy+point[1], col)
		}
		y++
		if y%2 == 0 {
			x--
		}

	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for i := range img.Pix {
		img.Pix[i] = 255 // white
	}
	line(img, 0, 0, 100, 100, RED)
	circle(img, 100, 100, 10, RED)
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}
