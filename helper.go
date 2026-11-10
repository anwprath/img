package main

import "image"

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
