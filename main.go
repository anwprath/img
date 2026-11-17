package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"

	"gonum.org/v1/gonum/mat"
)

var (
	width  = 800
	height = 800
)

// y-axis rotation
func rotate(arr [3]float64) [3]float64 {
	const piBy6 = math.Pi / 6
	cos, sin := math.Cos(piBy6), math.Cos(piBy6)

	return [3]float64{
		cos*arr[0] + sin*arr[2],
		arr[1],
		-sin*arr[0] + cos*arr[2],
	}
}

// viewer at (0,0,c)
func perspective(arr [3]float64) [3]float64 {
	scale := 1.0 - arr[2]/3.0

	return [3]float64{
		arr[0] / scale,
		arr[1] / scale,
		arr[2] / scale,
	}
}

func normalize(arr [3]float64) (int, int, int) {
	x := int((arr[0] + 1.0) * float64(width) / 2)
	y := int((1 + arr[1]) * float64(height) / 2)
	z := int((1 + arr[2]) * float64(255) / 2)
	return x, y, z
}

func Draw() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	model, err := loadOBJ("diablo3_pose.obj")
	if err != nil {
		panic("error loading .obj file: " + err.Error())
	}

	eye := [3]float64{-1, 0, 2}   // camera position
	center := [3]float64{0, 0, 0} // camera direction
	up := [3]float64{0, 1, 0}     // camera up vector

	setModelView(eye, center, up)                                                                         // build the ModelView   matrix
	setPerspectiveTransform(norm([3]float64{eye[0] - center[0], eye[1] - center[1], eye[2] - center[2]})) // build the Perspective matrix
	setViewportTransform(width/16, height/16, width*7/8, height*7/8)                                      // build the Viewport    matrix

	for i := range model.faces {
		var clip [3]mat.VecDense

		for _, j := range []int32{0, 1, 2} {
			v := model.vertices[model.faces[i][j]]

			vertex := mat.NewVecDense(4, []float64{v[0], v[1], v[2], 1})
			var temp mat.VecDense
			temp.MulVec(ModelView, vertex)

			clip[j].MulVec(Perspective, &temp)
		}

		rasterize(clip, zBuffer, img, color.RGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())})

	}
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
