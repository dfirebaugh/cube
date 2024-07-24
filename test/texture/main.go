package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/dfirebaugh/cube/pkg/block"
	"github.com/dfirebaugh/cube/pkg/primitive"
	"github.com/dfirebaugh/cube/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 600, "Textured Cube", nil, nil)
	if err != nil {
		log.Fatalln("failed to create glfw window:", err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize glow:", err)
	}

	gl.Enable(gl.DEPTH_TEST)

	program, err := shader.NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalf("failed to create program: %v", err)
	}
	gl.UseProgram(program)

	cube := block.TestBlock()
	textures := []uint32{cube.Front, cube.Back, cube.Left, cube.Right, cube.Top, cube.Bottom}
	for i, texture := range textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture)
	}

	for i := 0; i < len(textures); i++ {
		uniformName := fmt.Sprintf("textures[%d]", i)
		loc := gl.GetUniformLocation(program, gl.Str(uniformName+"\x00"))
		gl.Uniform1i(loc, int32(i))
	}

	start := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		angle := float32(time.Since(start).Seconds())
		model := mgl32.HomogRotate3DY(angle).Mul4(mgl32.HomogRotate3DX(angle * 0.5))

		modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		renderCube(&cube, program)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func renderCube(cube *primitive.Cube, program uint32) {
	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, 0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0,
		0.5, 0.5, -0.5, 1.0, 1.0, 0,
		0.5, 0.5, -0.5, 1.0, 1.0, 0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0,
		-0.5, -0.5, -0.5, 0.0, 0.0, 0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1,
		0.5, -0.5, 0.5, 1.0, 0.0, 1,
		0.5, 0.5, 0.5, 1.0, 1.0, 1,
		0.5, 0.5, 0.5, 1.0, 1.0, 1,
		-0.5, 0.5, 0.5, 0.0, 1.0, 1,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1,

		-0.5, 0.5, 0.5, 1.0, 0.0, 2,
		-0.5, 0.5, -0.5, 1.0, 1.0, 2,
		-0.5, -0.5, -0.5, 0.0, 1.0, 2,
		-0.5, -0.5, -0.5, 0.0, 1.0, 2,
		-0.5, -0.5, 0.5, 0.0, 0.0, 2,
		-0.5, 0.5, 0.5, 1.0, 0.0, 2,

		0.5, 0.5, 0.5, 1.0, 0.0, 3,
		0.5, 0.5, -0.5, 1.0, 1.0, 3,
		0.5, -0.5, -0.5, 0.0, 1.0, 3,
		0.5, -0.5, -0.5, 0.0, 1.0, 3,
		0.5, -0.5, 0.5, 0.0, 0.0, 3,
		0.5, 0.5, 0.5, 1.0, 0.0, 3,

		-0.5, -0.5, -0.5, 0.0, 1.0, 4,
		0.5, -0.5, -0.5, 1.0, 1.0, 4,
		0.5, -0.5, 0.5, 1.0, 0.0, 4,
		0.5, -0.5, 0.5, 1.0, 0.0, 4,
		-0.5, -0.5, 0.5, 0.0, 0.0, 4,
		-0.5, -0.5, -0.5, 0.0, 1.0, 4,

		-0.5, 0.5, -0.5, 0.0, 1.0, 5,
		0.5, 0.5, -0.5, 1.0, 1.0, 5,
		0.5, 0.5, 0.5, 1.0, 0.0, 5,
		0.5, 0.5, 0.5, 1.0, 0.0, 5,
		-0.5, 0.5, 0.5, 0.0, 0.0, 5,
		-0.5, 0.5, -0.5, 0.0, 1.0, 5,
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(2, 1, gl.FLOAT, false, 6*4, 5*4)
	gl.EnableVertexAttribArray(2)

	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	gl.BindVertexArray(0)
}

var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in float aFace;

uniform mat4 model;

out vec2 TexCoord;
out float Face;

void main() {
    gl_Position = model * vec4(aPos, 1.0);
    TexCoord = aTexCoord;
    Face = aFace;
}
` + "\x00"

var fragmentShaderSource = `
#version 330 core
out vec4 FragColor;

in vec2 TexCoord;
in float Face;

uniform sampler2D textures[6];

void main() {
    int faceIndex = int(Face);
    FragColor = texture(textures[faceIndex], TexCoord);
}
` + "\x00"
