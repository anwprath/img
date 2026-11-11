package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand/v2"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

var (
	width  = 4096
	height = 4096
)

func normalize(arr [3]float64) (int, int) {
	x := int((arr[0] + 1.0) * float64(width) / 2)
	y := int((1 + arr[1]) * float64(height) / 2)

	return x, y
}

func Draw() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	model, err := loadOBJ("diablo3_pose.obj")
	if err != nil {
		panic("error loading .obj file: " + err.Error())
	}

	var wg sync.WaitGroup
	for i := range model.faces {
		ax, ay := normalize(model.vertices[model.faces[i][0]])
		bx, by := normalize(model.vertices[model.faces[i][1]])
		cx, cy := normalize(model.vertices[model.faces[i][2]])
		// wg.Go(func() {
			triangle(img, color.RGBA{
				uint8(rand.IntN(255)), uint8(rand.IntN(255)), uint8(rand.IntN(255)), 255}, ax, ay, bx, by, cx, cy)
		// })
	}

	wg.Wait()
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
