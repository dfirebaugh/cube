package main

import (
	"log"
	"math/rand"

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

	for x := 0; x < chunkWidth; x++ {
		for z := 0; z < chunkLength; z++ {
			height := rand.Intn(chunkHeight/2) + chunkHeight/4
			for y := 0; y < height; y++ {
				color := getColorForHeight(y)
				meshRenderer.AddCube(primitive.Cube{
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

	e.Run()
}

func getColorForHeight(y int) component.Color {
	if y < chunkHeight/4 {
		return component.Color{0.6, 0.4, 0.2}
	} else if y < chunkHeight/2 {
		return component.Color{0.2, 1.0, 0.2}
	} else {
		return component.Color{0.8, 0.8, 0.8}
	}
}
