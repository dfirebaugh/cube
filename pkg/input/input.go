package input

import (
	"github.com/dfirebaugh/cube/pkg/message"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sirupsen/logrus"
)

var (
	isWireframeEnabled bool
	lastX              float64
	lastY              float64
	firstMouse         = true
	keyStates          = make(map[glfw.Key]glfw.Action)
)

func Update(window *glfw.Window, broker message.MessageBus) {
	handleMovement(window, broker)
	handleVerticalMovement(window, broker, 0.1)
	handleMouseMovement(window, broker)
	toggleWireframeMode(window, broker)

	updateKeyStates(window)
}

func handleMovement(window *glfw.Window, broker message.MessageBus) {
	var speed float32 = 0.05
	if window.GetKey(glfw.KeyW) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, 0, speed},
		})
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, 0, -speed},
		})
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{-speed, 0, 0},
		})
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{speed, 0, 0},
		})
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, speed, 0},
		})
	}
	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, -speed, 0},
		})
	}
}

func handleMouseMovement(window *glfw.Window, broker message.MessageBus) {
	x, y := window.GetCursorPos()

	if firstMouse {
		lastX = x
		lastY = y
		firstMouse = false
	}

	xOffset := float32(x - lastX)
	yOffset := float32(lastY - y) // Reversed since y-coordinates go from bottom to top
	lastX = x
	lastY = y

	broker.Publish(message.Message{
		Topic:     "MouseMovement",
		Requestor: "input",
		Payload:   [2]float32{xOffset, yOffset},
	})
}

func toggleWireframeMode(window *glfw.Window, broker message.MessageBus) {
	if !IsButtonJustPressed(window, glfw.KeyZ) {
		return
	}

	event := message.Message{
		Topic:     "ToggleWireframe",
		Requestor: "input",
		Payload:   nil,
	}
	broker.Publish(message.Message{
		Topic:     "TestEvent",
		Requestor: "input",
		Payload:   nil,
	})
	logrus.Infoln("Publishing ToggleWireframe event")
	broker.Publish(event)
}

func handleVerticalMovement(window *glfw.Window, broker message.MessageBus, zIncrement float32) {
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, zIncrement, 0},
		})
	}

	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		broker.Publish(message.Message{
			Topic:     "CameraMove",
			Requestor: "input",
			Payload:   [3]float32{0, -zIncrement, 0},
		})
	}
}

func updateKeyStates(window *glfw.Window) {
	for key := glfw.KeySpace; key <= glfw.KeyLast; key++ {
		keyStates[key] = window.GetKey(key)
	}
}

func IsButtonJustPressed(window *glfw.Window, key glfw.Key) bool {
	currentState := window.GetKey(key)
	if currentState == glfw.Press && keyStates[key] == glfw.Release {
		return true
	}
	return false
}
