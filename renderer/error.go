package renderer

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/sirupsen/logrus"
)

func checkGLError(location string) {
	for {
		err := gl.GetError()
		if err == 0 {
			break
		}
		logrus.Errorf("OpenGL error at %s: %d\n", location, err)
	}
}
