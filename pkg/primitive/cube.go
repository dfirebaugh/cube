package primitive

import (
	"github.com/dfirebaugh/cube/pkg/component"
)

type Cube struct {
	component.Position
	component.Color
	Size float32
	CubeTexture
	ShouldHide bool
	HideFront  bool
	HideBack   bool
	HideLeft   bool
	HideRight  bool
	HideTop    bool
	HideBottom bool
}

type CubeTexture struct {
	Front  uint32
	Back   uint32
	Left   uint32
	Right  uint32
	Top    uint32
	Bottom uint32
}

func (c Cube) Vertices() []float32 {
	halfSize := c.Size / 2
	var vertices []float32

	if c.ShouldHide {
		return vertices
	}

	if !c.HideFront {
		vertices = append(vertices, []float32{
			-halfSize, -halfSize, halfSize, 0.0, 0.0, 0.0, 0.0, 1.0,
			halfSize, -halfSize, halfSize, 1.0, 0.0, 0.0, 0.0, 1.0,
			halfSize, halfSize, halfSize, 1.0, 1.0, 0.0, 0.0, 1.0,
			-halfSize, halfSize, halfSize, 0.0, 1.0, 0.0, 0.0, 1.0,
		}...)
	}

	if !c.HideBack {
		vertices = append(vertices, []float32{
			-halfSize, -halfSize, -halfSize, 0.0, 0.0, 0.0, 0.0, -1.0,
			halfSize, -halfSize, -halfSize, 1.0, 0.0, 0.0, 0.0, -1.0,
			halfSize, halfSize, -halfSize, 1.0, 1.0, 0.0, 0.0, -1.0,
			-halfSize, halfSize, -halfSize, 0.0, 1.0, 0.0, 0.0, -1.0,
		}...)
	}

	if !c.HideLeft {
		vertices = append(vertices, []float32{
			-halfSize, -halfSize, -halfSize, 0.0, 0.0, -1.0, 0.0, 0.0,
			-halfSize, -halfSize, halfSize, 1.0, 0.0, -1.0, 0.0, 0.0,
			-halfSize, halfSize, halfSize, 1.0, 1.0, -1.0, 0.0, 0.0,
			-halfSize, halfSize, -halfSize, 0.0, 1.0, -1.0, 0.0, 0.0,
		}...)
	}

	if !c.HideRight {
		vertices = append(vertices, []float32{
			halfSize, -halfSize, -halfSize, 0.0, 0.0, 1.0, 0.0, 0.0,
			halfSize, -halfSize, halfSize, 1.0, 0.0, 1.0, 0.0, 0.0,
			halfSize, halfSize, halfSize, 1.0, 1.0, 1.0, 0.0, 0.0,
			halfSize, halfSize, -halfSize, 0.0, 1.0, 1.0, 0.0, 0.0,
		}...)
	}

	if !c.HideTop {
		vertices = append(vertices, []float32{
			-halfSize, halfSize, -halfSize, 0.0, 0.0, 0.0, 1.0, 0.0,
			halfSize, halfSize, -halfSize, 1.0, 0.0, 0.0, 1.0, 0.0,
			halfSize, halfSize, halfSize, 1.0, 1.0, 0.0, 1.0, 0.0,
			-halfSize, halfSize, halfSize, 0.0, 1.0, 0.0, 1.0, 0.0,
		}...)
	}

	if !c.HideBottom {
		vertices = append(vertices, []float32{
			-halfSize, -halfSize, -halfSize, 0.0, 0.0, 0.0, -1.0, 0.0,
			halfSize, -halfSize, -halfSize, 1.0, 0.0, 0.0, -1.0, 0.0,
			halfSize, -halfSize, halfSize, 1.0, 1.0, 0.0, -1.0, 0.0,
			-halfSize, -halfSize, halfSize, 0.0, 1.0, 0.0, -1.0, 0.0,
		}...)
	}

	return vertices
}

func (c Cube) Indices(offset uint32) []uint32 {
	var indices []uint32
	var currentOffset uint32 = 0

	if c.ShouldHide {
		return indices
	}
	if !c.HideFront {
		indices = append(indices, []uint32{
			currentOffset + 0, currentOffset + 1, currentOffset + 2,
			currentOffset + 2, currentOffset + 3, currentOffset + 0,
		}...)
		currentOffset += 4
	}

	if !c.HideBack {
		indices = append(indices, []uint32{
			currentOffset + 2, currentOffset + 1, currentOffset + 0,
			currentOffset + 0, currentOffset + 3, currentOffset + 2,
		}...)
		currentOffset += 4
	}

	if !c.HideLeft {
		indices = append(indices, []uint32{
			currentOffset + 0, currentOffset + 1, currentOffset + 2,
			currentOffset + 2, currentOffset + 3, currentOffset + 0,
		}...)
		currentOffset += 4
	}

	if !c.HideRight {
		indices = append(indices, []uint32{
			currentOffset + 2, currentOffset + 1, currentOffset + 0,
			currentOffset + 0, currentOffset + 3, currentOffset + 2,
		}...)
		currentOffset += 4
	}

	if !c.HideTop {
		indices = append(indices, []uint32{
			currentOffset + 2, currentOffset + 1, currentOffset + 0,
			currentOffset + 0, currentOffset + 3, currentOffset + 2,
		}...)
		currentOffset += 4
	}

	if !c.HideBottom {
		indices = append(indices, []uint32{
			currentOffset + 0, currentOffset + 1, currentOffset + 2,
			currentOffset + 2, currentOffset + 3, currentOffset + 0,
		}...)
		currentOffset += 4
	}

	return indices
}
