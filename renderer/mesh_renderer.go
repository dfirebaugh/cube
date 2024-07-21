package renderer

import (
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type MeshRenderer struct {
	program   uint32
	cubes     []primitive.Cube
	wireframe bool
	camera    Camera
	window    Window
	bus       message.MessageBus
	events    chan string
	mesher    Mesher
	meshDirty bool
}

func NewMeshRenderer(mesher Mesher) *MeshRenderer {
	renderer := &MeshRenderer{
		events:    make(chan string),
		mesher:    mesher,
		meshDirty: true,
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

func (r *MeshRenderer) SetCamera(camera Camera) {
	r.camera = camera
}

func (r *MeshRenderer) SetWindow(window Window) {
	r.window = window
}

func (r *MeshRenderer) SetMessageBus(m message.MessageBus) {
	r.bus = m
	go r.subscribeToEvents()
}

func (r *MeshRenderer) AddCube(cube primitive.Cube) {
	r.cubes = append(r.cubes, cube)
	r.meshDirty = true
}

func (r *MeshRenderer) ToggleWireframe() {
	r.wireframe = !r.wireframe
	if r.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	checkGLError("ToggleWireframe")
}

func (r *MeshRenderer) Render() {
	if r.meshDirty {
		r.mesher.CreateMesh(r.cubes)
		r.meshDirty = false
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.program)
	checkGLError("UseProgram")
	r.mesher.Bind()
	checkGLError("BindMesh")

	r.SetShaderUniforms()

	r.mesher.Draw()
	checkGLError("DrawMesh")

	r.mesher.Unbind()
	checkGLError("UnbindMesh")

	r.drainEvents()
}

func (r *MeshRenderer) SetShaderUniforms() {
	width, height := r.window.GetSize()
	view := r.camera.GetViewMatrix()
	projection := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)

	viewLoc := gl.GetUniformLocation(r.program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(r.program, gl.Str("projection\x00"))

	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])
	checkGLError("SetShaderUniforms")
}

func (r *MeshRenderer) drainEvents() {
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

func (r *MeshRenderer) subscribeToEvents() {
	if r.bus == nil {
		logrus.Println("MessageBus not set for MeshRenderer")
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
