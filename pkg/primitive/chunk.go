package primitive

import "github.com/go-gl/mathgl/mgl32"

const ChunkSize = 16

type Chunk struct {
	blocks   [ChunkSize][ChunkSize][ChunkSize]Cube
	position mgl32.Vec3
}

func NewChunk(position mgl32.Vec3) *Chunk {
	return &Chunk{position: position}
}

func (c *Chunk) SetBlock(x, y, z int, cube Cube) {
	if x >= 0 && x < ChunkSize && y >= 0 && y < ChunkSize && z >= 0 && z < ChunkSize {
		c.blocks[x][y][z] = cube
	}
}

func (c *Chunk) GetBlock(x, y, z int) Cube {
	if x >= 0 && x < ChunkSize && y >= 0 && y < ChunkSize && z >= 0 && z < ChunkSize {
		return c.blocks[x][y][z]
	}
	return Cube{}
}

func (c *Chunk) WorldPosition() mgl32.Vec3 {
	return c.position
}

func (c *Chunk) IsFaceExposed(x, y, z int, face string) bool {
	switch face {
	case "left":
		return x == 0 || c.blocks[x-1][y][z].Size == 0
	case "right":
		return x == ChunkSize-1 || c.blocks[x+1][y][z].Size == 0
	case "bottom":
		return y == 0 || c.blocks[x][y-1][z].Size == 0
	case "top":
		return y == ChunkSize-1 || c.blocks[x][y+1][z].Size == 0
	case "back":
		return z == 0 || c.blocks[x][y][z-1].Size == 0
	case "front":
		return z == ChunkSize-1 || c.blocks[x][y][z+1].Size == 0
	}
	return false
}
