package renderer

import (
	"fmt"
	"unsafe"

	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GreedyMesher struct {
	vao      uint32
	vbo      uint32
	ebo      uint32
	vertices []float32
	indices  []uint32
}

const (
	chunkSize = 10
)

func NewGreedyMesher() *GreedyMesher {
	return &GreedyMesher{}
}

func (m *GreedyMesher) CreateMesh(cubes []primitive.Cube) {
	m.vertices = nil
	m.indices = nil
	m.createGreedyMesh(cubes)
	m.setupBuffers()
}

func (m *GreedyMesher) Bind() {
	gl.BindVertexArray(m.vao)
}

func (m *GreedyMesher) Unbind() {
	gl.BindVertexArray(0)
}

func (m *GreedyMesher) Draw() {
	gl.DrawElements(gl.TRIANGLES, int32(len(m.indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (m *GreedyMesher) GetMesh() ([]float32, []uint32) {
	return m.vertices, m.indices
}

func (m *GreedyMesher) String() string {
	return fmt.Sprintf("Vertices: %v\nIndices: %v", m.vertices, m.indices)
}

func (m *GreedyMesher) createGreedyMesh(cubes []primitive.Cube) {
	var vertices []float32
	var indices []uint32

	solid := make([][][]bool, chunkSize)
	for i := range solid {
		solid[i] = make([][]bool, chunkSize)
		for j := range solid[i] {
			solid[i][j] = make([]bool, chunkSize)
		}
	}

	cubeColors := make(map[[3]int]mgl32.Vec3)
	for _, cube := range cubes {
		pos := [3]int{int(cube.X), int(cube.Y), int(cube.Z)}
		solid[pos[0]][pos[1]][pos[2]] = true
		cubeColors[pos] = cube.Color
	}

	for d := 0; d < 3; d++ {
		u := (d + 1) % 3
		v := (d + 2) % 3
		x := [3]int{0, 0, 0}
		q := [3]int{0, 0, 0}
		q[d] = 1

		mask := make([]bool, chunkSize*chunkSize)

		for x[d] = -1; x[d] < chunkSize; {
			n := 0
			for x[v] = 0; x[v] < chunkSize; x[v]++ {
				for x[u] = 0; x[u] < chunkSize; x[u]++ {
					blockCurrent := x[d] >= 0 && solid[x[0]][x[1]][x[2]]
					blockCompare := x[d] < chunkSize-1 && solid[x[0]+q[0]][x[1]+q[1]][x[2]+q[2]]

					mask[n] = blockCurrent != blockCompare
					n++
				}
			}

			x[d]++

			n = 0
			for j := 0; j < chunkSize; j++ {
				for i := 0; i < chunkSize; {
					if mask[n] {
						w := 1
						for w+i < chunkSize && mask[n+w] {
							w++
						}

						h := 1
						done := false
						for h+j < chunkSize {
							for k := 0; k < w; k++ {
								if !mask[n+k+h*chunkSize] {
									done = true
									break
								}
							}
							if done {
								break
							}
							h++
						}

						x[u], x[v] = i, j
						du := [3]int{0, 0, 0}
						dv := [3]int{0, 0, 0}
						du[u] = w
						dv[v] = h

						color := mgl32.Vec3{1, 0, 0} // Default color
						if x[d] >= 0 && x[d] < chunkSize && x[u] >= 0 && x[u] < chunkSize && x[v] >= 0 && x[v] < chunkSize {
							pos := [3]int{x[0], x[1], x[2]}
							if col, exists := cubeColors[pos]; exists {
								color = col
							}
						}

						vertices = append(vertices,
							float32(x[0]), float32(x[1]), float32(x[2]), color[0], color[1], color[2],
							float32(x[0]+du[0]), float32(x[1]+du[1]), float32(x[2]+du[2]), color[0], color[1], color[2],
							float32(x[0]+dv[0]), float32(x[1]+dv[1]), float32(x[2]+dv[2]), color[0], color[1], color[2],
							float32(x[0]+du[0]+dv[0]), float32(x[1]+du[1]+dv[1]), float32(x[2]+du[2]+dv[2]), color[0], color[1], color[2],
						)

						idx := uint32(len(vertices)/6 - 4)
						indices = append(indices,
							idx, idx+1, idx+2,
							idx+1, idx+3, idx+2,
						)

						for l := 0; l < h; l++ {
							for k := 0; k < w; k++ {
								mask[n+k+l*chunkSize] = false
							}
						}

						i += w
						n += w
					} else {
						i++
						n++
					}
				}
			}
		}
	}

	m.vertices = vertices
	m.indices = indices
}

func (m *GreedyMesher) setupBuffers() {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*4, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.indices)*4, gl.Ptr(m.indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	m.vao = vao
	m.vbo = vbo
	m.ebo = ebo
}
