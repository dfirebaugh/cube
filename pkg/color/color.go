package color

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
)

// RGBAtoVec3 converts color.RGBA to mgl32.Vec3
func RGBAtoVec3(c color.RGBA) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(c.R) / 255.0,
		float32(c.G) / 255.0,
		float32(c.B) / 255.0,
	}
}

// Vec3toRGBA converts mgl32.Vec3 to color.RGBA
func Vec3toRGBA(v mgl32.Vec3) color.RGBA {
	return color.RGBA{
		R: uint8(v[0] * 255.0),
		G: uint8(v[1] * 255.0),
		B: uint8(v[2] * 255.0),
		A: 255,
	}
}
