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
	cubeSize = 10
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

	for x := 0; x < cubeSize; x++ {
		for y := 0; y < cubeSize; y++ {
			for z := 0; z < cubeSize; z++ {
				cube := block.TestBlock()
				cube.X = float32(x)
				cube.Y = float32(y)
				cube.Z = float32(z)
				cube.Color = component.Color{float32(x) / float32(cubeSize), float32(y) / float32(cubeSize), float32(z) / float32(cubeSize)}
				cubeRenderer.AddCube(cube)
			}
		}
	}

	e.Run()
}
