package main

import (
	"image"
	"image/color"
)

func triangle(img *image.RGBA, c color.RGBA, x0, y0, x1, y1, x2, y2 int) {
	line(img, c , x0, y0, x1, y1)
	line(img, c , x0, y0, x2, y2)
	line(img, c , x1, y1, x2, y2)
}

// x = x0 + m*(x1-x0)
// y = y0 + m*(y1-y0)
// now replace m in second eqn with m=(x-x0)/(x1-x0) from first equation
// multiply everything by 2 to remove floating point math
// 63, 0, 32,32
func line(img *image.RGBA, c color.RGBA, x0, y0, x1, y1 int) {
	steep := abs(y1-y0) > abs(x1-x0)

	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	y := y0
	verticalDirection := 1

	if y1 < y0 {
		verticalDirection = -1
	}
	xDiff := x1 - x0
	deviation := 0
	for x := x0; x <= x1; x++ {
		var swX, swY int = x, y
		if steep {
			swX, swY = y, x
		}
		if (image.Point{swX, swY}.In(img.Rect)) {
			i := img.PixOffset(swX, swY)
			s := img.Pix[i : i+4 : i+4]
			s[0] = c.R
			s[1] = c.G
			s[2] = c.B
			s[3] = c.A
		}

		deviation += 2 * abs((y1 - y0))
		if deviation > xDiff {
			y += verticalDirection
			deviation -= 2 * (xDiff)
		}
	}
}
