package main

import (
	"log"

	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/component"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/renderer"
)

const (
	chunkWidth  = 16
	chunkLength = 16
	chunkHeight = 16
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

	// Generate cat shape
	generateCat(meshRenderer)

	e.Run()
}

// generateCat generates a simple cat shape using blocks
func generateCat(renderer *renderer.MeshRenderer) {
	// Body
	for x := 2; x < 6; x++ {
		for z := 2; z < 6; z++ {
			for y := 0; y < 4; y++ {
				color := component.Color{0.8, 0.8, 0.8} // Grey for the cat's body
				renderer.AddCube(primitive.Cube{
					Position: component.Position{
						X: float32(x),
						Y: float32(y),
						Z: float32(z),
					},
					Size:  1.0,
					Color: color,
				})
			}
		}
	}

	// Head
	for x := 3; x < 5; x++ {
		for z := 1; z < 3; z++ {
			for y := 4; y < 6; y++ {
				color := component.Color{0.8, 0.8, 0.8} // Grey for the cat's head
				renderer.AddCube(primitive.Cube{
					Position: component.Position{
						X: float32(x),
						Y: float32(y),
						Z: float32(z),
					},
					Size:  1.0,
					Color: color,
				})
			}
		}
	}

	// Ears
	for x := 3; x < 5; x++ {
		z := 1
		y := 6
		color := component.Color{0.8, 0.8, 0.8} // Grey for the cat's ears
		renderer.AddCube(primitive.Cube{
			Position: component.Position{
				X: float32(x),
				Y: float32(y),
				Z: float32(z),
			},
			Size:  1.0,
			Color: color,
		})
	}

	// Legs
	for x := 2; x < 6; x += 3 {
		for z := 2; z < 6; z += 3 {
			for y := 0; y < 2; y++ {
				color := component.Color{0.8, 0.8, 0.8} // Grey for the cat's legs
				renderer.AddCube(primitive.Cube{
					Position: component.Position{
						X: float32(x),
						Y: float32(y),
						Z: float32(z),
					},
					Size:  1.0,
					Color: color,
				})
			}
		}
	}

	// Tail
	for x := 6; x < 8; x++ {
		z := 4
		y := 3
		color := component.Color{0.8, 0.8, 0.8} // Grey for the cat's tail
		renderer.AddCube(primitive.Cube{
			Position: component.Position{
				X: float32(x),
				Y: float32(y),
				Z: float32(z),
			},
			Size:  1.0,
			Color: color,
		})
	}
}
