package engine

import (
	"log"
	"runtime"

	"github.com/dfirebaugh/cube/pkg/camera"
	"github.com/dfirebaugh/cube/pkg/input"
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/dfirebaugh/cube/pkg/message/broker"
	"github.com/dfirebaugh/cube/pkg/player"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/sirupsen/logrus"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	player    *player.Player
	window    *glfw.Window
	camera    *camera.Camera
	renderers []renderer.Renderer
	bus       message.MessageBus
}

var worldHasLoaded bool

func init() {
	runtime.LockOSThread()
}

func New(startup func()) *Engine {
	go func() {
		startup()
		worldHasLoaded = true
	}()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	monitor := glfw.GetPrimaryMonitor()
	videoMode := monitor.GetVideoMode()

	windowWidth := 800
	windowHeight := 450
	windowX := (videoMode.Width - windowWidth) / 2
	windowY := (videoMode.Height - windowHeight) / 2

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "cube", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.SetPos(windowX, windowY)
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	engine := &Engine{
		player: player.New(-5, 0, 15, broker.NewBroker()),
		window: window,
		camera: camera.NewCamera(window),
		bus:    broker.NewBroker(),
	}

	input.Init(window, engine.bus)

	window.SetFramebufferSizeCallback(engine.framebufferSizeCallback)

	engine.camera.ListenForInputEvents(engine.bus)
	engine.camera.UpdatePosition(-2, 0, 10)

	go engine.bus.Start()
	go engine.subscribeToEvents()

	return engine
}

func (e *Engine) framebufferSizeCallback(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func (e *Engine) close() {
	e.bus.Stop()
	glfw.Terminate()
}

func (e *Engine) AddRenderer(renderer renderer.Renderer) {
	renderer.SetCamera(e.camera)
	renderer.SetWindow(e.window)
	renderer.SetMessageBus(e.bus)
	e.renderers = append(e.renderers, renderer)
}

func (e *Engine) applyPhysics() {
}

func (e *Engine) update() {
	input.Update(e.window, e.bus)

	if !worldHasLoaded {
		return
	}

	e.applyPhysics()
}

func (e *Engine) ShouldClose() bool {
	return e.window.ShouldClose()
}

func (e *Engine) SwapBuffers() {
	e.window.SwapBuffers()
	glfw.PollEvents()
}

func (e Engine) ClearScreen() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
}

func (e *Engine) draw() {
	for _, r := range e.renderers {
		r.Render()
	}
}

func (e *Engine) Run() {
	defer e.close()
	for !e.ShouldClose() {
		e.update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		e.ClearScreen()
		e.draw()

		e.SwapBuffers()
	}
}

func (e *Engine) subscribeToEvents() {
	if e.bus == nil {
		logrus.Println("MessageBus not set for CubeRenderer")
		return
	}

	msg := e.bus.Subscribe()
	defer e.bus.Unsubscribe(msg)

	for {
		select {
		case m, ok := <-msg:
			if !ok {
				return
			}
			if m.GetTopic() == "RequestClose" {
				e.window.SetShouldClose(true)
			}
		}
	}
}
