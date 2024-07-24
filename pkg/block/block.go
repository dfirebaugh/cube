package block

import (
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/sirupsen/logrus"
)

var (
	red    uint32
	green  uint32
	blue   uint32
	yellow uint32
	grey   uint32
	pink   uint32

	texturesLoaded bool
)

func loadTextures() {
	var err error
	red, err = renderer.LoadTexture("assets/textures/red.png")
	if err != nil {
		logrus.Fatalln("failed to load back texture:", err)
	}
	green, err = renderer.LoadTexture("assets/textures/green.png")
	if err != nil {
		logrus.Fatalln("failed to load front texture:", err)
	}
	blue, err = renderer.LoadTexture("assets/textures/blue.png")
	if err != nil {
		logrus.Fatalln("failed to load left texture:", err)
	}
	yellow, err = renderer.LoadTexture("assets/textures/yellow.png")
	if err != nil {
		logrus.Fatalln("failed to load right texture:", err)
	}
	grey, err = renderer.LoadTexture("assets/textures/grey.png")
	if err != nil {
		logrus.Fatalln("failed to load top texture:", err)
	}
	pink, err = renderer.LoadTexture("assets/textures/pink.png")
	if err != nil {
		logrus.Fatalln("failed to load bottom texture:", err)
	}

	texturesLoaded = true
}

func TestBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  red,
			Back:   green,
			Left:   blue,
			Right:  yellow,
			Top:    grey,
			Bottom: pink,
		},
	}
}

func PinkBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  pink,
			Back:   pink,
			Left:   pink,
			Right:  pink,
			Top:    pink,
			Bottom: pink,
		},
	}
}

func RedBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  red,
			Back:   red,
			Left:   red,
			Right:  red,
			Top:    red,
			Bottom: red,
		},
	}
}

func BlueBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  blue,
			Back:   blue,
			Left:   blue,
			Right:  blue,
			Top:    blue,
			Bottom: blue,
		},
	}
}

func GreenBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  green,
			Back:   green,
			Left:   green,
			Right:  green,
			Top:    green,
			Bottom: green,
		},
	}
}

func YellowBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  yellow,
			Back:   yellow,
			Left:   yellow,
			Right:  yellow,
			Top:    yellow,
			Bottom: yellow,
		},
	}
}

func GreyBlock() primitive.Cube {
	if !texturesLoaded {
		loadTextures()
	}
	return primitive.Cube{
		Size: 1.0,
		CubeTexture: primitive.CubeTexture{
			Front:  grey,
			Back:   grey,
			Left:   grey,
			Right:  grey,
			Top:    grey,
			Bottom: grey,
		},
	}
}
