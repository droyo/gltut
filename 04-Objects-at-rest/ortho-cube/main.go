// ortho-cube displays a prism in 3D space without perspective projection.
// It is an implementation of http://arcsynthesis.org/gltut/Positioning/Tut04%20Perspective%20Projection.html
package main

import (
	"log"
	"time"
	"aqwari.net/exp/gl"
	"aqwari.net/exp/display"
)

var config = display.Config{
	"Title":          "Perspective Projection",
	"Geometry":       "500x500",
	"OpenGL Version": "3.2",
}

var vertShader = []byte(
`#version 150

in vec4 position;
in vec4 color;

smooth out vec4 theColor;
uniform vec2 offset;

void main() {
	gl_Position = position + vec4(offset, 0, 0);
	theColor = color;
}
`)

var fragShader = []byte(
`#version 150

smooth in vec4 theColor;
out vec4 outColor;

void main() {
	outColor = theColor;
}
`)

func main() {
	win, err := display.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Close()
	if err := gl.Init(config["OpenGL Version"]); err != nil {
		log.Fatal(err)
	}
	
	gl.ClearColor(0, 0, 0, 0)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	
	vertexData := []float32{
		0.25, 0.25, 0.75, 1.0,
		0.25, -0.25, 0.75, 1.0,
		-0.25, 0.25, 0.75, 1.0,

		0.25, -0.25, 0.75, 1.0,
		-0.25, -0.25, 0.75, 1.0,
		-0.25, 0.25, 0.75, 1.0,

		0.25, 0.25, -0.75, 1.0,
		-0.25, 0.25, -0.75, 1.0,
		0.25, -0.25, -0.75, 1.0,

		0.25, -0.25, -0.75, 1.0,
		-0.25, 0.25, -0.75, 1.0,
		-0.25, -0.25, -0.75, 1.0,

		-0.25, 0.25, 0.75, 1.0,
		-0.25, -0.25, 0.75, 1.0,
		-0.25, -0.25, -0.75, 1.0,

		-0.25, 0.25, 0.75, 1.0,
		-0.25, -0.25, -0.75, 1.0,
		-0.25, 0.25, -0.75, 1.0,

		0.25, 0.25, 0.75, 1.0,
		0.25, -0.25, -0.75, 1.0,
		0.25, -0.25, 0.75, 1.0,

		0.25, 0.25, 0.75, 1.0,
		0.25, 0.25, -0.75, 1.0,
		0.25, -0.25, -0.75, 1.0,

		0.25, 0.25, -0.75, 1.0,
		0.25, 0.25, 0.75, 1.0,
		-0.25, 0.25, 0.75, 1.0,

		0.25, 0.25, -0.75, 1.0,
		-0.25, 0.25, 0.75, 1.0,
		-0.25, 0.25, -0.75, 1.0,

		0.25, -0.25, -0.75, 1.0,
		-0.25, -0.25, 0.75, 1.0,
		0.25, -0.25, 0.75, 1.0,

		0.25, -0.25, -0.75, 1.0,
		-0.25, -0.25, -0.75, 1.0,
		-0.25, -0.25, 0.75, 1.0,

		0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 1.0, 1.0,

		0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 1.0, 1.0,

		0.8, 0.8, 0.8, 1.0,
		0.8, 0.8, 0.8, 1.0,
		0.8, 0.8, 0.8, 1.0,

		0.8, 0.8, 0.8, 1.0,
		0.8, 0.8, 0.8, 1.0,
		0.8, 0.8, 0.8, 1.0,

		0.0, 1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,

		0.0, 1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,

		0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, 0.0, 1.0,

		0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, 0.0, 1.0,

		1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 0.0, 1.0,

		1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 0.0, 1.0,

		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,

		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
		0.0, 1.0, 1.0, 1.0,
	}
	
	prog := gl.CreateProgram()
	defer gl.DeleteProgram(prog)
	
	vert := gl.CreateShader(gl.VERTEX_SHADER)
	defer gl.DeleteShader(vert)
	
	frag := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(frag)
	
	gl.ShaderSource(vert, vertShader)
	gl.ShaderSource(frag, fragShader)
	
	if err := gl.CompileShader(vert); err != nil {
		log.Fatal(err)
	}
	if err := gl.CompileShader(frag); err != nil {
		log.Fatal(err)
	}
	
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	if err := gl.LinkProgram(prog); err != nil {
		log.Fatal(err)
	}
	gl.DetachShader(prog, vert)
	gl.DetachShader(prog, frag)
	
	gl.UseProgram(prog)
	
	buffers := gl.GenBuffers(1)
	defer gl.DeleteBuffers(buffers)
	
	gl.BindBuffer(gl.ARRAY_BUFFER, buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, vertexData, gl.STATIC_DRAW)
	
	arr := gl.GenVertexArrays(1)
	gl.BindVertexArray(arr[0])
	
	pos, _ := gl.GetAttribLocation(prog, "position")
	col, _ := gl.GetAttribLocation(prog, "color")
	gl.EnableVertexAttribArray(pos)
	gl.EnableVertexAttribArray(col)
	
	offset, _ := gl.GetUniformLocation(prog, "offset")
	gl.Uniformf(offset, 0.5, 0.25)
	
	gl.VertexAttribPointer(pos, 4, gl.Float32, false, 0, 0)
	gl.VertexAttribPointer(col, 4, gl.Float32, false, 0, uintptr(len(vertexData))/2 * 4)
	
	clock := time.Tick(time.Second / 30)
	gl.Clear(gl.COLOR_BUFFER_BIT)
Loop:
	for _ = range clock {
		select {
		case ev := <-win.Event:
			switch ev := ev.(type) {
			case display.KeyPress:
				if ev.Code == display.KeyEscape {
					break Loop
				}
			case display.Resize:
				gl.Viewport(0, 0, ev.Width, ev.Height)
			}
		default:
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
			win.Flip()
			win.CheckEvent()
		}
	}
}
