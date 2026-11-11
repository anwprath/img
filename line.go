package main

import (
	"image"
	"image/color"
)

// x = x0 + m*(x1-x0)
// y = y0 + m*(y1-y0)
// now replace m in second eqn with m=(x-x0)/(x1-x0) from first equation
// multiply everything by 2 to remove floating point math
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

	y := int(y0)
	verticalDirection := 1

	if y1 < y0 {
		verticalDirection = -1
	}
	xDiff := x1 - x0
	deviation := 0
	for x := x0; x <= x1; x++ {
		// When steep, coordinates were swapped before the loop, so the
		// pixel must be transposed back into real image space.
		px, py := x, y
		if steep {
			px, py = y, x
		}
		if (image.Point{px, py}).In(img.Rect) {
			i := img.PixOffset(px, py)
			s := img.Pix[i : i+4 : i+4]
			s[0] = c.R
			s[1] = c.G
			s[2] = c.B
			s[3] = c.A
		}

		// Advance the error term on every step, including pixels that fall
		// outside the image, so clipping does not skew the rest of the line.
		deviation += 2 * abs(y1-y0)
		if deviation > xDiff {
			y += verticalDirection
			deviation -= 2 * xDiff
		}
	}
}
