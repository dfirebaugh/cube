package main

import (
	"log"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/block"
	"github.com/dfirebaugh/cube/pkg/component"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/sirupsen/logrus"
)

const (
	cubeSize = 40
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	e := engine.New(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in startup function:", r)
			}
		}()
	})

	cubeRenderer := renderer.NewMeshRenderer(renderer.NewCubeMesher())
	e.AddRenderer(cubeRenderer)
	cube := block.TestBlock()
	cube.Size = 1
	cube.Color = component.Color{float32(0) / float32(cubeSize), float32(255) / float32(cubeSize), float32(0) / float32(cubeSize)}
	cube.Position = component.Position{
		X: float32(1),
		Y: 0,
		Z: 0,
	}

	cubeRenderer.AddCube(cube)

	e.Run()
}
