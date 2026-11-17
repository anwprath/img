package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gonum.org/v1/gonum/mat"
)

func rasterize(clip [3]mat.VecDense, zbuffer map[xyPoint]float64, img *image.RGBA, color color.RGBA) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	// normalized device coordinates [-1,1]
	var ndc [3][3]float64
	var screen [3][2]float64

	for i := 0; i < 3; i++ {
		w := clip[i].AtVec(3)
		if w == 0 {
			return
		}

		ndc[i] = [3]float64{
			clip[i].AtVec(0) / w,
			clip[i].AtVec(1) / w,
			clip[i].AtVec(2) / w,
		}

		var viewportResult mat.VecDense
		viewportInput := mat.NewVecDense(4, []float64{
			ndc[i][0],
			ndc[i][1],
			ndc[i][2],
			1,
		})
		viewportResult.MulVec(Viewport, viewportInput)

		screen[i] = [2]float64{
			viewportResult.AtVec(0),
			viewportResult.AtVec(1),
		}
	}

	abc := mat.NewDense(3, 3, []float64{
		screen[0][0], screen[0][1], 1,
		screen[1][0], screen[1][1], 1,
		screen[2][0], screen[2][1], 1,
	})

	// Backface culling and discard subpixel triangles. (area<1)
	if mat.Det(abc) < 1 {
		return
	}

	minX := math.Min(screen[0][0], math.Min(screen[1][0], screen[2][0]))
	maxX := math.Max(screen[0][0], math.Max(screen[1][0], screen[2][0]))
	minY := math.Min(screen[0][1], math.Min(screen[1][1], screen[2][1]))
	maxY := math.Max(screen[0][1], math.Max(screen[1][1], screen[2][1]))

	bbMinX := max(0, int(math.Floor(minX)))
	bbMaxX := min(width-1, int(math.Ceil(maxX)))
	bbMinY := max(0, int(math.Floor(minY)))
	bbMaxY := min(height-1, int(math.Ceil(maxY)))

	var abcInveT mat.Dense
	if err := abcInveT.Inverse(abc.T()); err != nil {
		fmt.Println("error inverting abc: ", err)
		return
	}

	for x := bbMinX; x <= bbMaxX; x++ {
		for y := bbMinY; y <= bbMaxY; y++ {
			point := mat.NewVecDense(3, []float64{
				float64(x),
				float64(y),
				1,
			})

			var bc mat.VecDense
			bc.MulVec(&abcInveT, point)

			alpha := bc.AtVec(0)
			beta := bc.AtVec(1)
			gamma := bc.AtVec(2)

			// Point outside triangle.
			if alpha < 0 || beta < 0 || gamma < 0 {
				continue
			}

			// Interpolate depth.
			z := alpha*ndc[0][2] +
				beta*ndc[1][2] +
				gamma*ndc[2][2]

			p := xyPoint{x: x, y: y}

			// Untouched pixels (missing map entry, reads as 0) must not
			// reject fragments with negative NDC depth.
			if zCurrent, ok := zbuffer[p]; ok && z <= zCurrent {
				continue
			}

			zbuffer[p] = z
			img.SetRGBA(x, y, color)
		}
	}
}

// AI generated code below since I don't know how to use gonum
var (
	Viewport    = mat.NewDense(4, 4, nil)
	Perspective = mat.NewDense(4, 4, nil)
	ModelView   = mat.NewDense(4, 4, nil)
)

func setViewportTransform(x, y, w, h int) {
	Viewport = mat.NewDense(4, 4, []float64{
		float64(w) / 2, 0, 0, float64(x) + float64(w)/2,
		0, float64(h) / 2, 0, float64(y) + float64(h)/2,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
}

func setPerspectiveTransform(f float64) {
	Perspective = mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, -1 / f, 1,
	})
}

func setModelView(eye, center, up [3]float64) {
	n := normalized(sub(eye, center))
	l := normalized(cross(up, n))
	m := normalized(cross(n, l))

	rotation := mat.NewDense(4, 4, []float64{
		l[0], l[1], l[2], 0,
		m[0], m[1], m[2], 0,
		n[0], n[1], n[2], 0,
		0, 0, 0, 1,
	})

	translation := mat.NewDense(4, 4, []float64{
		1, 0, 0, -center[0],
		0, 1, 0, -center[1],
		0, 0, 1, -center[2],
		0, 0, 0, 1,
	})

	ModelView = mat.NewDense(4, 4, nil)
	ModelView.Mul(rotation, translation)
}

func normalized(v [3]float64) [3]float64 {
	magnitude := math.Sqrt(
		v[0]*v[0] +
			v[1]*v[1] +
			v[2]*v[2],
	)

	if magnitude == 0 {
		return [3]float64{}
	}

	return [3]float64{
		v[0] / magnitude,
		v[1] / magnitude,
		v[2] / magnitude,
	}
}

func norm(v [3]float64) float64 {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

func cross(a, b [3]float64) [3]float64 {
	return [3]float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func sub(a, b [3]float64) [3]float64 {
	return [3]float64{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}
