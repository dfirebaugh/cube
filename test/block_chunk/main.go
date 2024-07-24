package main

import (
	"math/rand"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/block"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	e := engine.New(func() {
	})

	cubeRenderer := renderer.NewBlockRenderer()
	e.AddRenderer(cubeRenderer)

	chunk := primitive.NewChunk()

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				var cube primitive.Cube
				switch rand.Intn(6) {
				case 0:
					cube = block.TestBlock()
				case 1:
					cube = block.RedBlock()
				case 2:
					cube = block.BlueBlock()
				case 3:
					cube = block.PinkBlock()
				case 4:
					cube = block.GreenBlock()
				case 5:
					cube = block.YellowBlock()
				case 6:
					cube = block.GreyBlock()
				}

				chunk.SetBlock(x, y, z, cube)
			}
		}
	}

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				cube := chunk.GetBlock(x, y, z)
				if cube.Size > 0 {
					cube.X = float32(x)
					cube.Y = float32(y)
					cube.Z = float32(z)

					// Hide faces that are not exposed
					cube.HideLeft = !chunk.IsFaceExposed(x, y, z, "left")
					cube.HideRight = !chunk.IsFaceExposed(x, y, z, "right")
					cube.HideBottom = !chunk.IsFaceExposed(x, y, z, "bottom")
					cube.HideTop = !chunk.IsFaceExposed(x, y, z, "top")
					cube.HideBack = !chunk.IsFaceExposed(x, y, z, "back")
					cube.HideFront = !chunk.IsFaceExposed(x, y, z, "front")

					cubeRenderer.AddCube(cube)
				}
			}
		}
	}

	e.Run()
}
