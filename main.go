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
