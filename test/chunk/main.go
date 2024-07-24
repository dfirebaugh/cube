package main

import (
	"math/rand"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/block"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	e := engine.New(func() {
	})

	chunkRenderer := renderer.NewChunkRenderer(mgl32.Vec3{0, 0, 0})
	e.AddRenderer(chunkRenderer)

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

				chunkRenderer.SetBlock(x, y, z, cube)
			}
		}
	}

	e.Run()
}
