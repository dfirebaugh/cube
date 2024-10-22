package main

import (
	"log"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/component"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
)

const (
	cubeSize = 10
)

func main() {
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
					Position: component.Position{
						X: float32(x),
						Y: float32(y),
						Z: float32(z),
					},
					Size:  1.0,
					Color: component.Color{float32(255), 0, 0},
				})
			}
		}
	}
	e.Run()
}
