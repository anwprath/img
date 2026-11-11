package main

import (
	"image"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
)

func Draw() {
	width := 64
	height := 64
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	_, err := loadOBJ("diablo3_pose.obj")
	if err != nil {
		panic("error loading .obj file: " + err.Error())
	}

	line(img, red, 0, 0, 32, 32)
	line(img, blue, 0, 64, 32, 32)
	line(img, white, 64, 64, 32, 32)
	line(img, green, 64, 0, 32, 32)

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
