package renderer

import (
	"fmt"
	"unsafe"

	"github.com/dfirebaugh/cube/pkg/component"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/go-gl/gl/v3.3-core/gl"
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
	m.generateMesh(cubes)
	m.setupBuffers()
}

func (m *GreedyMesher) Bind() {
	gl.BindVertexArray(m.vao)
}

func (m *GreedyMesher) Unbind() {
	gl.BindVertexArray(0)
}

func (m *GreedyMesher) Draw() {
	m.EnableBackfaceCulling()
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.BindVertexArray(0)
}

func (m *GreedyMesher) GetMesh() ([]float32, []uint32) {
	return m.vertices, m.indices
}

func (m *GreedyMesher) String() string {
	return fmt.Sprintf("Vertices: %v\nIndices: %v", m.vertices, m.indices)
}

func (m *GreedyMesher) populateSolidAndColors(cubes []primitive.Cube) ([][][]bool, map[[3]int]component.Color) {
	const expandedChunkSize = 15
	solid := make([][][]bool, expandedChunkSize)
	for i := range solid {
		solid[i] = make([][]bool, expandedChunkSize)
		for j := range solid[i] {
			solid[i][j] = make([]bool, expandedChunkSize)
		}
	}

	cubeColors := make(map[[3]int]component.Color)
	for _, cube := range cubes {
		pos := [3]int{int(cube.X), int(cube.Y), int(cube.Z)}
		if pos[0] >= 0 && pos[0] < expandedChunkSize && pos[1] >= 0 && pos[1] < expandedChunkSize && pos[2] >= 0 && pos[2] < expandedChunkSize {
			solid[pos[0]][pos[1]][pos[2]] = true
			cubeColors[pos] = cube.Color
		}
	}
	return solid, cubeColors
}

func (m *GreedyMesher) generateMesh(cubes []primitive.Cube) {
	const expandedChunkSize = 15
	solid, cubeColors := m.populateSolidAndColors(cubes)
	for d := 0; d < 3; d++ {
		m.generateDirectionMesh(d, solid, cubeColors, expandedChunkSize)
	}
}

func (m *GreedyMesher) generateDirectionMesh(d int, solid [][][]bool, cubeColors map[[3]int]component.Color, expandedChunkSize int) {
	u := (d + 1) % 3
	v := (d + 2) % 3
	x := [3]int{0, 0, 0}
	q := [3]int{0, 0, 0}
	q[d] = 1

	visited := make([]bool, (expandedChunkSize+1)*(expandedChunkSize+1))

	for x[d] = -1; x[d] < expandedChunkSize; {
		n := 0
		for x[v] = 0; x[v] < expandedChunkSize; x[v]++ {
			for x[u] = 0; x[u] < expandedChunkSize; x[u]++ {
				currentBlock := m.isNotEmptyBlock(x, d, solid)
				compareBlock := m.isNotEmptyBlock([3]int{x[0] + q[0], x[1] + q[1], x[2] + q[2]}, d, solid)

				visited[n] = currentBlock != compareBlock
				n++
			}
		}

		x[d]++

		n = 0
		for j := 0; j < expandedChunkSize; j++ {
			for i := 0; i < expandedChunkSize; {
				if visited[n] {
					w, h := m.findWidthAndHeight(visited, expandedChunkSize, n, i, j)

					x[u], x[v] = i, j
					du := [3]int{0, 0, 0}
					dv := [3]int{0, 0, 0}
					du[u] = w
					dv[v] = h

					color := m.getColor(x, d, expandedChunkSize, cubeColors)

					if m.isNotEmptyBlock(x, d, solid) {
						m.generatePositiveFace(d, x, du, dv, color)
					} else {
						m.generateNegativeFace(d, x, du, dv, color)
					}

					m.markAsVisited(visited, expandedChunkSize, n, w, h)

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

func (m *GreedyMesher) isNotEmptyBlock(x [3]int, d int, solid [][][]bool) bool {
	return x[d] >= 0 && x[0] < len(solid) && x[1] < len(solid) && x[2] < len(solid) && solid[x[0]][x[1]][x[2]]
}

func (m *GreedyMesher) findWidthAndHeight(visited []bool, expandedChunkSize, n, i, j int) (int, int) {
	w := 1
	for w+i < expandedChunkSize && visited[n+w] {
		w++
	}

	h := 1
	done := false
	for h+j < expandedChunkSize {
		for k := 0; k < w; k++ {
			if !visited[n+k+h*expandedChunkSize] {
				done = true
				break
			}
		}
		if done {
			break
		}
		h++
	}
	return w, h
}

func (m *GreedyMesher) getColor(x [3]int, d, expandedChunkSize int, cubeColors map[[3]int]component.Color) component.Color {
	color := component.Color{1, 0, 0} // Default color
	if x[d] >= 0 && x[d] < expandedChunkSize && x[0] >= 0 && x[0] < expandedChunkSize && x[1] >= 0 && x[1] < expandedChunkSize && x[2] >= 0 && x[2] < expandedChunkSize {
		pos := [3]int{x[0], x[1], x[2]}
		if col, exists := cubeColors[pos]; exists {
			color = col
		}
	}
	return color
}

func (m *GreedyMesher) generatePositiveFace(d int, x, du, dv [3]int, color component.Color) {
	if d == 0 { // Positive X face
		m.addFaceVertices(x, du, dv, color)
		m.addPositiveXFaceIndices()
	} else if d == 1 { // Positive Y face
		m.addFaceVertices(x, du, dv, color)
		m.addPositiveYFaceIndices()
	} else if d == 2 { // Positive Z face
		m.addFaceVertices(x, du, dv, color)
		m.addPositiveZFaceIndices()
	}
}

func (m *GreedyMesher) generateNegativeFace(d int, x, du, dv [3]int, color component.Color) {
	if d == 0 { // Negative X face
		m.addFaceVertices(x, dv, du, color)
		m.addNegativeXFaceIndices()
	} else if d == 1 { // Negative Y face
		m.addFaceVertices(x, dv, du, color)
		m.addNegativeYFaceIndices()
	} else if d == 2 { // Negative Z face
		m.addFaceVertices(x, dv, du, color)
		m.addNegativeZFaceIndices()
	}
}

func (m *GreedyMesher) addFaceVertices(x, du, dv [3]int, color component.Color) {
	m.vertices = append(m.vertices,
		float32(x[0]), float32(x[1]), float32(x[2]), color[0], color[1], color[2],
		float32(x[0]+du[0]), float32(x[1]+du[1]), float32(x[2]+du[2]), color[0], color[1], color[2],
		float32(x[0]+dv[0]), float32(x[1]+dv[1]), float32(x[2]+dv[2]), color[0], color[1], color[2],
		float32(x[0]+du[0]+dv[0]), float32(x[1]+du[1]+dv[1]), float32(x[2]+du[2]+dv[2]), color[0], color[1], color[2],
	)
}

func (m *GreedyMesher) addPositiveXFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx+2, idx+1, idx, // First triangle
		idx+2, idx+3, idx+1, // Second triangle
	)
}

func (m *GreedyMesher) addPositiveYFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx, idx+2, idx+1, // First triangle
		idx+1, idx+2, idx+3, // Second triangle
	)
}

func (m *GreedyMesher) addPositiveZFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx+2, idx+1, idx, // First triangle
		idx+2, idx+3, idx+1, // Second triangle
	)
}

func (m *GreedyMesher) addNegativeXFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx, idx+2, idx+1, // First triangle
		idx+1, idx+2, idx+3, // Second triangle
	)
}

func (m *GreedyMesher) addNegativeYFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx+2, idx+1, idx, // First triangle
		idx+2, idx+3, idx+1, // Second triangle
	)
}

func (m *GreedyMesher) addNegativeZFaceIndices() {
	idx := uint32(len(m.vertices)/6 - 4)
	m.indices = append(m.indices,
		idx, idx+2, idx+1, // First triangle
		idx+1, idx+2, idx+3, // Second triangle
	)
}

func (m *GreedyMesher) markAsVisited(visited []bool, expandedChunkSize, n, w, h int) {
	for l := 0; l < h; l++ {
		for k := 0; k < w; k++ {
			visited[n+k+l*expandedChunkSize] = false
		}
	}
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

func (m *GreedyMesher) EnableBackfaceCulling() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
}
