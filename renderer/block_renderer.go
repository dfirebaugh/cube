package renderer

import (
	"fmt"

	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type BlockRenderer struct {
	camera      Camera
	window      Window
	messageBus  message.MessageBus
	cubeProgram uint32
	cubes       []primitive.Cube
	wireframe   bool
	events      chan string
}

func (r *BlockRenderer) SetCamera(camera Camera) {
	r.camera = camera
}

func (r *BlockRenderer) SetWindow(window Window) {
	r.window = window
}

func (r *BlockRenderer) SetMessageBus(bus message.MessageBus) {
	r.messageBus = bus
	go r.subscribeToEvents()
}

func (r *BlockRenderer) ToggleWireframe() {
	r.wireframe = !r.wireframe
	if r.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	checkGLError("ToggleWireframe")
}

func (r *BlockRenderer) Render() {
	if r.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	view := r.camera.GetViewMatrix()

	windowWidth, windowHeight := r.window.GetSize()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0)

	viewLoc := gl.GetUniformLocation(r.cubeProgram, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(r.cubeProgram, gl.Str("projection\x00"))

	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	for _, cube := range r.cubes {
		model := mgl32.Translate3D(cube.Position.X, cube.Position.Y, cube.Position.Z)
		modelLoc := gl.GetUniformLocation(r.cubeProgram, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

		r.renderCube(cube)
	}
	r.drainEvents()
}

func (r *BlockRenderer) renderCube(cube primitive.Cube) {
	vertices := cube.Vertices()
	indices := cube.Indices(0)

	if len(vertices) == 0 || len(indices) == 0 {
		return
	}

	textures := []uint32{cube.Front, cube.Back, cube.Left, cube.Right, cube.Top, cube.Bottom}
	for i, texture := range textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}

	for i := 0; i < len(textures); i++ {
		uniformName := fmt.Sprintf("textures[%d]", i)
		loc := gl.GetUniformLocation(r.cubeProgram, gl.Str(uniformName+"\x00"))
		gl.Uniform1i(loc, int32(i))
	}

	var VAO, VBO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 8*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 8*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, 8*4, 5*4)
	gl.EnableVertexAttribArray(2)

	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteBuffers(1, &EBO)
}

func (r *BlockRenderer) AddCube(cube primitive.Cube) {
	r.cubes = append(r.cubes, cube)
}

func (r *BlockRenderer) drainEvents() {
	for {
		select {
		case event := <-r.events:
			if event == "ToggleWireframe" {
				r.ToggleWireframe()
				logrus.Trace("Received ToggleWireframe event")
			}
		default:
			return // Exit the loop when there are no more events
		}
	}
}

func (r *BlockRenderer) subscribeToEvents() {
	if r.messageBus == nil {
		logrus.Println("MessageBus not set for BlockRenderer")
		return
	}

	msg := r.messageBus.Subscribe()
	defer r.messageBus.Unsubscribe(msg)

	for {
		select {
		case m, ok := <-msg:
			if !ok {
				return
			}
			if m.GetTopic() == "ToggleWireframe" {
				r.events <- m.GetTopic()
			}
		default:
		}
	}
}

func NewBlockRenderer() *BlockRenderer {
	vertexShaderSource, err := shader.ShaderFS.ReadFile("block_vertex_shader.glsl")
	if err != nil {
		logrus.Fatalf("failed to read vertex shader: %v", err)
	}

	fragmentShaderSource, err := shader.ShaderFS.ReadFile("block_fragment_shader.glsl")
	if err != nil {
		logrus.Fatalf("failed to read fragment shader: %v", err)
	}
	cubeProgram, err := shader.NewProgram(string(vertexShaderSource)+"\x00", string(fragmentShaderSource)+"\x00")
	if err != nil {
		logrus.Fatalf("failed to create program: %v", err)
	}
	gl.UseProgram(cubeProgram)

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.DEPTH_TEST)

	return &BlockRenderer{
		cubeProgram: cubeProgram,
		cubes:       []primitive.Cube{},
		events:      make(chan string),
	}
}
