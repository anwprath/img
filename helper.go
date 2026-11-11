package main

import "image"

type Numeric interface {
	int | float64 | int32
}

func abs[T Numeric](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func flipVertical(img *image.RGBA) {
	b := img.Bounds()

	for y := 0; y < b.Dy()/2; y++ {
		oppositeY := b.Dy() - 1 - y

		for x := 0; x < b.Dx(); x++ {
			top := img.RGBAAt(x, y)
			bottom := img.RGBAAt(x, oppositeY)

			img.SetRGBA(x, y, bottom)
			img.SetRGBA(x, oppositeY, top)
		}
	}
}
