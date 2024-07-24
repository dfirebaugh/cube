package renderer

import (
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	GetPosition() mgl32.Vec3
	GetDirection() mgl32.Vec3
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
	CreateMesh(cubes []primitive.Cube)
	Bind()
	Unbind()
	Draw()
	GetMesh() ([]float32, []uint32)
	String() string
}
