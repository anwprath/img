package main

import (
	"image"
	"image/png"
	"os"
	"runtime"
	"runtime/pprof"
)

const width = 800
const height = 800

func project(arr [3]float64) (int, int) {
	return int((arr[0] + 1.0) * (float64(width) / 2)), int((arr[1] + 1.0) * (float64(height) / 2))
}

func Draw() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	obj, err := loadOBJ("diablo3_pose.obj")
	if err != nil {
		panic("error loading .obj file: " + err.Error())
	}

	for i := range obj.faces {
		ax, ay := project(obj.faceVertex(i, 0))
		bx, by := project(obj.faceVertex(i, 1))
		cx, cy := project(obj.faceVertex(i, 2))
		line(img, red, ax, ay, bx, by)
		line(img, red, bx, by, cx, cy)
		line(img, red, cx, cy, ax, ay)
	}

	for i := 0; i < obj.vertexCount(); i++ {
		v := obj.vertex(i) // get i-th vertex
		x, y := project(v) // project it to the screen
		img.Set(x, y, white)

	}
	flipVertical(img) // inverted coordinate system of Go who tf knows why

	f1, _ := os.Create("out.png")
	defer f1.Close()
	png.Encode(f1, img)
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
