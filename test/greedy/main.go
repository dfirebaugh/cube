package main

import (
	"log"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	cubeSize = 10
)

func main() {
	log.Println("World created")

	e := engine.New(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in startup function:", r)
			}
		}()
	})

	meshRenderer := renderer.NewMeshRenderer(renderer.NewGreedyMesher())
	e.AddRenderer(meshRenderer)

	for x := 0; x < cubeSize; x++ {
		for y := 0; y < cubeSize; y++ {
			for z := 0; z < cubeSize; z++ {
				meshRenderer.AddCube(primitive.Cube{
					X:    float32(x),
					Y:    float32(y),
					Z:    float32(z),
					Size: 1.0,
					// Color: mgl32.Vec3{float32(x) / float32(cubeSize), float32(y) / float32(cubeSize), float32(z) / float32(cubeSize)},
					Color: mgl32.Vec3{float32(255), 0, 0},
				})
			}
		}
	}
	e.Run()
}
