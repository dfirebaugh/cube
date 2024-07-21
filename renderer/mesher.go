package renderer

import (
	"fmt"

	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type CubeMesher struct {
	vao      uint32
	vbo      uint32
	vertices []float32
}

func NewCubeMesher() *CubeMesher {
	return &CubeMesher{}
}

func (m *CubeMesher) CreateMesh(cubes []primitive.Cube) {
	m.vertices = nil
	m.createCube(cubes)
	m.setupBuffers()
}

func (m *CubeMesher) Bind() {
	gl.BindVertexArray(m.vao)
}

func (m *CubeMesher) Unbind() {
	gl.BindVertexArray(0)
}

func (m *CubeMesher) Draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.vertices)/6))
}

func (m *CubeMesher) GetMesh() ([]float32, []uint32) {
	return m.vertices, nil
}

func (m *CubeMesher) String() string {
	return fmt.Sprintf("Vertices: %v", m.vertices)
}

func (m *CubeMesher) createCube(cubes []primitive.Cube) {
	for _, cube := range cubes {
		color := cube.Color
		// Cube vertices positions and colors
		cubeVertices := []float32{
			// Front face
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],

			// Back face
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],

			// Left face
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],

			// Right face
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],

			// Top face
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y + cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],

			// Bottom face
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X + cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z + cube.Size/2, color[0], color[1], color[2],
			cube.X - cube.Size/2, cube.Y - cube.Size/2, cube.Z - cube.Size/2, color[0], color[1], color[2],
		}

		m.vertices = append(m.vertices, cubeVertices...)
	}
}

func (m *CubeMesher) setupBuffers() {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*4, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	m.vao = vao
	m.vbo = vbo
}
