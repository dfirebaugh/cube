package camera

import (
	"math"

	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

var Mode CameraMode

type CameraMode int

const (
	Custom CameraMode = iota
	FirstPerson
	Orbit
)

func init() {
	Mode = FirstPerson
}

type Camera struct {
	window    *glfw.Window
	position  mgl32.Vec3
	direction mgl32.Vec3
	up        mgl32.Vec3
	right     mgl32.Vec3
	fov       float32
	Yaw       float32
	Pitch     float32
	Distance  float32
}

func NewCamera(window *glfw.Window) *Camera {
	return &Camera{
		window:    window,
		fov:       45.0,
		position:  mgl32.Vec3{0, 0, 0},
		direction: mgl32.Vec3{0, 0, 0},
		up:        mgl32.Vec3{0, 1, 0},
		right:     mgl32.Vec3{1, 0, 0},
		Yaw:       -90.0, // Set initial yaw to -90.0 to look forward
		Pitch:     0.0,   // Set initial pitch to 0.0 to look level
		Distance:  10.0,  // Set initial distance for orbit mode
	}
}

func (c *Camera) GetPosition() mgl32.Vec3 {
	return c.position
}

func (c *Camera) SetPosition(x, y, z float32) {
	c.position = mgl32.Vec3{x, y, z}
}

func (c *Camera) GetDirection() mgl32.Vec3 {
	return c.direction
}

func (c *Camera) SetDirection(direction mgl32.Vec3) {
	c.direction = direction
}

func (c *Camera) UpdatePosition(dx, dy, dz float32) {
	c.position = c.position.Add(mgl32.Vec3{dx, dy, dz})
}

func (c *Camera) Move(dx, dy, dz float32) {
	if Mode == FirstPerson {
		c.position = c.position.Add(c.direction.Mul(dz))
		c.position = c.position.Add(c.right.Mul(dx))
		c.position = c.position.Add(c.up.Mul(dy))
	} else {
		c.UpdatePosition(dx, dy, dz)
	}
}

func (c *Camera) updateOrbit(x, y, z float32) {
	c.position = mgl32.Vec3{
		x + c.Distance*float32(math.Cos(float64(c.Yaw)*math.Pi/180.0)*math.Cos(float64(c.Pitch)*math.Pi/180.0)),
		y + c.Distance*float32(math.Sin(float64(c.Pitch)*math.Pi/180.0)),
		z + c.Distance*float32(math.Sin(float64(c.Yaw)*math.Pi/180.0)*math.Cos(float64(c.Pitch)*math.Pi/180.0)),
	}
	c.direction = mgl32.Vec3{
		-x, -y, -z,
	}.Sub(c.position).Normalize()
	c.updateCameraVectors()
}

func (c *Camera) Position() mgl32.Vec3 {
	return c.position
}

func (c *Camera) Direction() mgl32.Vec3 {
	return c.direction
}

func (c *Camera) Right() mgl32.Vec3 {
	return c.right
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.position, c.position.Add(c.direction), c.up)
}

func (c *Camera) GetProjectionMatrix() mgl32.Mat4 {
	windowWidth, windowHeight := c.window.GetSize()
	aspectRatio := float32(windowWidth) / float32(windowHeight)
	logrus.Infof("Window size: (%d, %d), Aspect ratio: %f", windowWidth, windowHeight, aspectRatio)
	return mgl32.Perspective(mgl32.DegToRad(c.fov), aspectRatio, 0.1, 100.0)
}

func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32, constrainPitch bool) {
	sensitivity := float32(0.1)
	xOffset *= sensitivity
	yOffset *= sensitivity

	c.Yaw += xOffset
	c.Pitch += yOffset

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	c.updateCameraVectors()
}

func (c *Camera) updateCameraVectors() {
	front := mgl32.Vec3{
		float32(math.Cos(float64(c.Yaw)*math.Pi/180.0) * math.Cos(float64(c.Pitch)*math.Pi/180.0)),
		float32(math.Sin(float64(c.Pitch) * math.Pi / 180.0)),
		float32(math.Sin(float64(c.Yaw)*math.Pi/180.0) * math.Cos(float64(c.Pitch)*math.Pi/180.0)),
	}
	c.direction = front.Normalize()
	c.right = c.direction.Cross(mgl32.Vec3{0, 1, 0}).Normalize()
	c.up = c.right.Cross(c.direction).Normalize()
}

func (c *Camera) ListenForInputEvents(broker message.MessageBus) {
	sub := broker.Subscribe()
	go func() {
		for msg := range sub {
			switch msg.GetTopic() {
			case "CameraMove":
				payload := msg.GetPayload().([3]float32)
				c.Move(payload[0], payload[1], payload[2])
			case "MouseMovement":
				payload := msg.GetPayload().([2]float32)
				c.ProcessMouseMovement(payload[0], payload[1], true)
			}
		}
	}()
}
