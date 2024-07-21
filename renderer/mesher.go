package renderer

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type CubeMesher struct {
	vao uint32
	vbo uint32
	ebo uint32
}

func NewCubeMesher() *CubeMesher {
	m := &CubeMesher{}
	m.createCube()
	return m
}

func (m *CubeMesher) CreateMesh() {
	m.createCube()
}

func (m *CubeMesher) Bind() {
	gl.BindVertexArray(m.vao)
}

func (m *CubeMesher) Unbind() {
	gl.BindVertexArray(0)
}

func (m *CubeMesher) Draw() {
	gl.DrawElements(gl.TRIANGLES, 36, gl.UNSIGNED_INT, unsafe.Pointer(nil))
}

func (m *CubeMesher) createCube() {
	vertices := []float32{
		// Positions        // Colors
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, // 0
		0.5, -0.5, -0.5, 0.0, 1.0, 0.0, // 1
		0.5, 0.5, -0.5, 0.0, 0.0, 1.0, // 2
		-0.5, 0.5, -0.5, 1.0, 1.0, 0.0, // 3
		-0.5, -0.5, 0.5, 1.0, 0.0, 1.0, // 4
		0.5, -0.5, 0.5, 0.0, 1.0, 1.0, // 5
		0.5, 0.5, 0.5, 1.0, 1.0, 1.0, // 6
		-0.5, 0.5, 0.5, 0.0, 0.0, 0.0, // 7
	}

	indices := []uint32{
		0, 1, 2, 2, 3, 0, // Back face
		4, 5, 6, 6, 7, 4, // Front face
		0, 1, 5, 5, 4, 0, // Bottom face
		2, 3, 7, 7, 6, 2, // Top face
		0, 3, 7, 7, 4, 0, // Left face
		1, 2, 6, 6, 5, 1, // Right face
	}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	m.vao = vao
	m.vbo = vbo
	m.ebo = ebo
}
