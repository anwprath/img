package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
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
// now replace m in second eqn with m=(x-x0)/(x1-x0) from first equation
// multiply everything by 2 to remove floating point math
func line(img *image.RGBA, c color.RGBA, x0, y0, x1, y1 int) {
	slope := abs(float64(y1-y0) / float64(x1-x0))

	if slope > 1 {
		line(img, c, y0, x0, y1, x1)
		return
	}

	if x0 > x1 {
		line(img, c, x1, y1, x0, y0)
		return
	}

	y := int(y0)
	verticalDirection := 1

	if y1 < y0 {
		verticalDirection = -1
	}
	xDiff := x1 - x0
	deviation := 0
	for x := x0; x <= x1; x++ {
		img.Set(x, int(y), c)
		deviation += 2 * abs((y1 - y0))
		if deviation > xDiff {
			y += verticalDirection
			deviation -= 2 * (xDiff)
		}
	}
}

func Draw() {
	width := 64
	height := 64
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	ax, ay := 7, 3
	bx, by := 12, 37
	cx, cy := 62, 53

	for i := 0; i < (1 << 24); i++ {
		ax, ay := rand.Intn(width), rand.Intn(height)
		bx, by := rand.Intn(width), rand.Intn(height)
		color := color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: uint8(rand.Intn(255)),
		}
		line(img, color, ax, ay, bx, by)

	}

	img.Set(ax, ay, white)
	img.Set(bx, by, white)
	img.Set(cx, cy, white)

	// flipVertical(img) // inverted coordinate system of Go who tf knows why
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	// -------- CPU profiling ----------------------------
	cpuProf, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(cpuProf); err != nil {
		panic(err)
	}
	defer cpuProf.Close()
	defer pprof.StopCPUProfile()

	// -------- Memory profiling -------------------------
	memProf, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memProf.Close()

	Draw()

	runtime.GC() // Force GC before collecting the heap profile.

	if err := pprof.WriteHeapProfile(memProf); err != nil {
		panic(err)
	}

}
