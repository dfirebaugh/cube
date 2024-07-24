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
	log.Println("World created")

	e := engine.New(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in startup function:", r)
			}
		}()
	})

	meshRenderer := renderer.NewMeshRenderer(renderer.NewGreedyMesher())
	// meshRenderer := renderer.NewMeshRenderer(renderer.NewCubeMesher())
	e.AddRenderer(meshRenderer)

	// Add the main block of cubes
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
					Color: component.Color{float32(x) / float32(cubeSize), float32(y) / float32(cubeSize), float32(z) / float32(cubeSize)},
					// Color: mgl32.Vec3{255, 0, 0},
				})
			}
		}
	}

	// Add a few extra cubes emerging from the main block in various positions
	extraCubes := []primitive.Cube{
		{Position: component.Position{X: 5, Y: 10, Z: 5}, Size: 1.0, Color: component.Color{0, 255, 0}},
		{Position: component.Position{X: 5, Y: 11, Z: 5}, Size: 1.0, Color: component.Color{0, 255, 0}},
		{Position: component.Position{X: 5, Y: 12, Z: 5}, Size: 1.0, Color: component.Color{0, 255, 0}},
		{Position: component.Position{X: 8, Y: 5, Z: 10}, Size: 1.0, Color: component.Color{0, 0, 255}},
		{Position: component.Position{X: 9, Y: 5, Z: 10}, Size: 1.0, Color: component.Color{0, 0, 255}},
		{Position: component.Position{X: 10, Y: 5, Z: 10}, Size: 1.0, Color: component.Color{0, 0, 255}},
		{Position: component.Position{X: 2, Y: 10, Z: 2}, Size: 1.0, Color: component.Color{255, 255, 0}},
		{Position: component.Position{X: 2, Y: 11, Z: 2}, Size: 1.0, Color: component.Color{255, 255, 0}},
		{Position: component.Position{X: 2, Y: 12, Z: 2}, Size: 1.0, Color: component.Color{255, 255, 0}},
		{Position: component.Position{X: 7, Y: 10, Z: 7}, Size: 1.0, Color: component.Color{0, 255, 255}},
		{Position: component.Position{X: 7, Y: 11, Z: 7}, Size: 1.0, Color: component.Color{0, 255, 255}},
		{Position: component.Position{X: 7, Y: 12, Z: 7}, Size: 1.0, Color: component.Color{0, 255, 255}},
		// Adding more cubes in varied and connected positions
		{Position: component.Position{X: 5, Y: 13, Z: 5}, Size: 1.0, Color: component.Color{128, 0, 128}},
		{Position: component.Position{X: 4, Y: 12, Z: 5}, Size: 1.0, Color: component.Color{128, 0, 128}},
		{Position: component.Position{X: 3, Y: 11, Z: 5}, Size: 1.0, Color: component.Color{128, 0, 128}},
		{Position: component.Position{X: 6, Y: 10, Z: 4}, Size: 1.0, Color: component.Color{255, 105, 180}},
		{Position: component.Position{X: 6, Y: 11, Z: 4}, Size: 1.0, Color: component.Color{255, 105, 180}},
		{Position: component.Position{X: 6, Y: 12, Z: 4}, Size: 1.0, Color: component.Color{255, 105, 180}},
		{Position: component.Position{X: 8, Y: 13, Z: 8}, Size: 1.0, Color: component.Color{75, 0, 130}},
		{Position: component.Position{X: 8, Y: 14, Z: 8}, Size: 1.0, Color: component.Color{75, 0, 130}},
		{Position: component.Position{X: 8, Y: 15, Z: 8}, Size: 1.0, Color: component.Color{75, 0, 130}},
	}

	for _, cube := range extraCubes {
		meshRenderer.AddCube(cube)
	}

	e.Run()
}
