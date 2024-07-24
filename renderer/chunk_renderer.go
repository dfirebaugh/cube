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

type ChunkRenderer struct {
	camera      Camera
	window      Window
	messageBus  message.MessageBus
	cubeProgram uint32
	chunk       *primitive.Chunk
	wireframe   bool
	events      chan string
}

func (r *ChunkRenderer) SetCamera(camera Camera) {
	r.camera = camera
}

func (r *ChunkRenderer) SetWindow(window Window) {
	r.window = window
}

func (r *ChunkRenderer) SetMessageBus(bus message.MessageBus) {
	r.messageBus = bus
	go r.subscribeToEvents()
}

func (r *ChunkRenderer) ToggleWireframe() {
	r.wireframe = !r.wireframe
	if r.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	checkGLError("ToggleWireframe")
}

func (r *ChunkRenderer) Render() {
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

	chunkPos := r.chunk.WorldPosition()
	for x := 0; x < primitive.ChunkSize; x++ {
		for y := 0; y < primitive.ChunkSize; y++ {
			for z := 0; z < primitive.ChunkSize; z++ {
				cube := r.chunk.GetBlock(x, y, z)
				if cube.Size > 0 {
					cube.Position.X = chunkPos.X() + float32(x)
					cube.Position.Y = chunkPos.Y() + float32(y)
					cube.Position.Z = chunkPos.Z() + float32(z)

					cube.HideLeft = !r.chunk.IsFaceExposed(x, y, z, "left")
					cube.HideRight = !r.chunk.IsFaceExposed(x, y, z, "right")
					cube.HideBottom = !r.chunk.IsFaceExposed(x, y, z, "bottom")
					cube.HideTop = !r.chunk.IsFaceExposed(x, y, z, "top")
					cube.HideBack = !r.chunk.IsFaceExposed(x, y, z, "back")
					cube.HideFront = !r.chunk.IsFaceExposed(x, y, z, "front")

					model := mgl32.Translate3D(cube.Position.X, cube.Position.Y, cube.Position.Z)
					modelLoc := gl.GetUniformLocation(r.cubeProgram, gl.Str("model\x00"))
					gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])

					r.renderCube(cube)
				}
			}
		}
	}
	r.drainEvents()
}

func (r *ChunkRenderer) renderCube(cube primitive.Cube) {
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

func (r *ChunkRenderer) SetBlock(x, y, z int, cube primitive.Cube) {
	r.chunk.SetBlock(x, y, z, cube)
}

func (r *ChunkRenderer) drainEvents() {
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

func (r *ChunkRenderer) subscribeToEvents() {
	if r.messageBus == nil {
		logrus.Println("MessageBus not set for ChunkRenderer")
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

func NewChunkRenderer(position mgl32.Vec3) *ChunkRenderer {
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

	return &ChunkRenderer{
		cubeProgram: cubeProgram,
		chunk:       primitive.NewChunk(position),
		events:      make(chan string),
	}
}
