package main

import (
	"image"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	width  = 150
	height = 150
)

func normalize(arr [3]float64) (int, int) {
	x := int((arr[0] + 1.0) * float64(width) / 2)
	y := int((1 + arr[1]) * float64(height) / 2)

	return x, y
}

func Draw() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	triangle(img, red, 7, 45, 35, 100, 45, 60)
	triangle(img, white, 120, 35, 90, 5, 45, 110)
	triangle(img, green, 115, 83, 80, 90, 85, 120)

	flipVertical(img) // inverted coordinate system of Go who tf knows why
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
