package renderer

import (
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type CubeRenderer struct {
	program   uint32
	cubes     []primitive.Cube
	wireframe bool
	camera    Camera
	window    Window
	bus       message.MessageBus
	events    chan string
	mesher    Mesher
}

func NewCubeRenderer(mesher Mesher) *CubeRenderer {
	renderer := &CubeRenderer{
		events: make(chan string),
		mesher: mesher,
	}
	vertexShaderSource, err := shader.ShaderFS.ReadFile("cube_vertex_shader.glsl")
	if err != nil {
		logrus.Fatalf("failed to read vertex shader: %v", err)
	}

	fragmentShaderSource, err := shader.ShaderFS.ReadFile("cube_fragment_shader.glsl")
	if err != nil {
		logrus.Fatalf("failed to read fragment shader: %v", err)
	}
	program, err := shader.NewProgram(string(vertexShaderSource)+"\x00", string(fragmentShaderSource)+"\x00")
	if err != nil {
		logrus.Fatalln("failed to create shader program:", err)
	}
	renderer.program = program

	return renderer
}

func (r *CubeRenderer) SetCamera(camera Camera) {
	r.camera = camera
}

func (r *CubeRenderer) SetWindow(window Window) {
	r.window = window
}

func (r *CubeRenderer) SetMessageBus(m message.MessageBus) {
	r.bus = m
	go r.subscribeToEvents()
}

func (r *CubeRenderer) AddCube(cube primitive.Cube) {
	r.cubes = append(r.cubes, cube)
}

func (r *CubeRenderer) ToggleWireframe() {
	r.wireframe = !r.wireframe
	if r.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	checkGLError("ToggleWireframe")
}

func (r *CubeRenderer) Render() {
	gl.Enable(gl.DEPTH_TEST)
	gl.UseProgram(r.program)
	r.mesher.Bind()

	for _, cube := range r.cubes {
		r.DrawCube(cube.X, cube.Y, cube.Z, cube.Size, cube.Color)
	}

	r.mesher.Unbind()

	// Drain events channel
	for {
		select {
		case event := <-r.events:
			if event == "ToggleWireframe" {
				r.ToggleWireframe()
			}
		default:
			return // Exit the loop when there are no more events
		}
	}
}

func (r *CubeRenderer) DrawCube(x, y, z, size float32, color mgl32.Vec3) {
	r.mesher.Bind()

	model := mgl32.Translate3D(x, y, z).Mul4(mgl32.Scale3D(size, size, size))
	width, height := r.window.GetSize()
	view := r.camera.GetViewMatrix()
	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)

	modelLoc := gl.GetUniformLocation(r.program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(r.program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(r.program, gl.Str("projection\x00"))
	colorLoc := gl.GetUniformLocation(r.program, gl.Str("inputColor\x00"))

	gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])
	gl.Uniform3fv(colorLoc, 1, &color[0])

	r.mesher.Draw()
	r.mesher.Unbind()
}

func (r *CubeRenderer) subscribeToEvents() {
	if r.bus == nil {
		logrus.Println("MessageBus not set for CubeRenderer")
		return
	}

	msg := r.bus.Subscribe()
	defer r.bus.Unsubscribe(msg)

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

func checkGLError(location string) {
	if err := gl.GetError(); err != 0 {
		logrus.Errorf("OpenGL error at %s: %d\n", location, err)
	}
}
