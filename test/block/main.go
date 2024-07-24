package main

import (
	"github.com/dfirebaugh/cube/engine"
	"github.com/dfirebaugh/cube/pkg/block"
	"github.com/dfirebaugh/cube/renderer"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	e := engine.New(func() {
		// Startup logic if needed
	})

	cubeRenderer := renderer.NewBlockRenderer()

	e.AddRenderer(cubeRenderer)

	test := block.TestBlock()
	test.X = 0
	test.Z = 0
	// test.ShouldHide = true
	// test.HideBack = true
	// test.HideTop = false
	// test.HideBottom = false
	// test.HideRight = true
	// test.HideLeft = true
	// test.HideFront = true
	red := block.RedBlock()
	red.X = 1
	red.Z = 1

	blue := block.BlueBlock()
	blue.X = 2
	blue.Z = 2

	cubeRenderer.AddCube(test)
	cubeRenderer.AddCube(red)
	cubeRenderer.AddCube(blue)

	e.Run()
}
