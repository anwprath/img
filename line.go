package main

import (
	"image"
	"image/color"
)

type Point struct {
	X, Y int
}

func signedTriangleArea(a, b, c Point) float64 {
	return 0.5 * float64((b.Y-a.Y)*(b.X+a.X)+(c.Y-b.Y)*(c.X+b.X)+(a.Y-c.Y)*(a.X+c.X))
}

func triangleBarycentricCoordinates(p, a, b, c Point) (float64, float64, float64) {
	totalArea := signedTriangleArea(a, b, c)
	alpha := signedTriangleArea(p, a, b) / totalArea
	beta := signedTriangleArea(p, b, c) / totalArea
	gamma := signedTriangleArea(p, c, a) / totalArea

	return alpha, beta, gamma
}

func pointInTriangle(a, b, c float64) bool {
	return !(a < 0 || b < 0 || c < 0)
}

func triangle(img *image.RGBA, c color.RGBA, x0, y0, x1, y1, x2, y2 int) {
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y0 > y2 {
		x0, y0, x2, y2 = x2, y2, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	line(img, c, x0, y0, x1, y1)
	line(img, c, x0, y0, x2, y2)
	line(img, c, x1, y1, x2, y2)

	for i := min(x0, x1, x2); i < max(x0, x1, x2); i++ {
		for j := min(y0, y1, y2); j < max(y0, y1, y2); j++ {
			alpha, beta, gamma := triangleBarycentricCoordinates(Point{i, j}, Point{x0, y0}, Point{x1, y1}, Point{x2, y2})
			if pointInTriangle(alpha, beta, gamma) {
				img.Set(i, j, c)
			}
		}
	}
}

// x = x0 + m*(x1-x0)
// y = y0 + m*(y1-y0)
// now replace m in second eqn with m=(x-x0)/(x1-x0) from first equation
// multiply everything b.Y 2 to remove floating point math
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
