package renderer

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/sirupsen/logrus"
)

func checkGLError(tag string) {
	for {
		err := gl.GetError()
		if err == gl.NO_ERROR {
			break
		}
		var errMsg string
		switch err {
		case gl.INVALID_ENUM:
			errMsg = "INVALID_ENUM"
		case gl.INVALID_VALUE:
			errMsg = "INVALID_VALUE"
		case gl.INVALID_OPERATION:
			errMsg = "INVALID_OPERATION"
		case gl.STACK_OVERFLOW:
			errMsg = "STACK_OVERFLOW"
		case gl.STACK_UNDERFLOW:
			errMsg = "STACK_UNDERFLOW"
		case gl.OUT_OF_MEMORY:
			errMsg = "OUT_OF_MEMORY"
		default:
			errMsg = "UNKNOWN"
		}
		logrus.Errorf("OpenGL error at %s: %d (%s)\n", tag, err, errMsg)
	}
}
