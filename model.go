package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Model is a triangle mesh loaded from a Wavefront .obj file.
type Model struct {
	// vertices holds the (x, y, z) position of every vertex.
	vertices [][3]float64
	// faces holds triangles as zero-based indices into vertices.
	faces [][3]int
}

// loadOBJ reads the Wavefront .obj file at path and returns its mesh.
// Only vertex positions and face definitions are needed by the renderer;
// texture coordinates, normals, and other directives are ignored.
// Polygonal faces are triangulated as a fan.
func loadOBJ(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	model := &Model{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "v":
			if len(fields) < 4 {
				return nil, fmt.Errorf("%s:%d: vertex needs 3 coordinates", path, lineNo)
			}

			var vertex [3]float64
			for i := range vertex {
				value, err := strconv.ParseFloat(fields[i+1], 64)
				if err != nil {
					return nil, fmt.Errorf("%s:%d: invalid vertex coordinate: %w", path, lineNo, err)
				}
				vertex[i] = value
			}
			model.vertices = append(model.vertices, vertex)

		case "f":
			if len(fields) < 4 {
				return nil, fmt.Errorf("%s:%d: face needs at least 3 vertices", path, lineNo)
			}

			indices := make([]int, 0, len(fields)-1)
			for _, field := range fields[1:] {
				index, err := parseVertexIndex(field, len(model.vertices))
				if err != nil {
					return nil, fmt.Errorf("%s:%d: %w", path, lineNo, err)
				}
				indices = append(indices, index)
			}

			// Triangulate the face as a fan.
			for i := 1; i+1 < len(indices); i++ {
				model.faces = append(model.faces, [3]int{indices[0], indices[i], indices[i+1]})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(model.vertices) == 0 {
		return nil, fmt.Errorf("%s: no vertices found", path)
	}

	return model, nil
}

// parseVertexIndex converts a face token such as "12", "12/3", or
// "12/3/4" into a zero-based vertex index. Negative OBJ indices are
// resolved relative to the vertices parsed so far.
func parseVertexIndex(token string, vertexCount int) (int, error) {
	if i := strings.IndexByte(token, '/'); i >= 0 {
		token = token[:i]
	}

	index, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("invalid face index %q", token)
	}
	if index < 0 {
		index = vertexCount + index
	} else {
		index--
	}
	if index < 0 || index >= vertexCount {
		return 0, fmt.Errorf("face index %q out of range", token)
	}

	return index, nil
}
