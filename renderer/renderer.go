package renderer

import (
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	GetViewMatrix() mgl32.Mat4
}

type Window interface {
	GetSize() (int, int)
}

type Renderer interface {
	Render()
	SetCamera(camera Camera)
	SetWindow(window Window)
	SetMessageBus(m message.MessageBus)
}

type Mesher interface {
	CreateMesh()
	Bind()
	Unbind()
	Draw()
}
