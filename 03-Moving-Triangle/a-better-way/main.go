package main

import (
	"log"
	"time"
	"math"
	"aqwari.net/exp/gl"
	"aqwari.net/exp/display"
)

var config = display.Config{
	"Title":          "A Better Way",
	"Geometry":       "500x500",
	"OpenGL Version": "3.2",
}

var vertShader = []byte(
`#version 150

in vec2 position;
uniform vec2 offset;

void main() {
	gl_Position = vec4(position + offset, 0, 1);
}
`)

var fragShader = []byte(
`#version 150

out vec4 outColor;
void main() {
	outColor = vec4(1,1,1,1);
}`)

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
	
	vertexData := []float32 {
		 0.0,    0.25,
		 0.25, -0.366,
		-0.25, -0.366,
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
	gl.EnableVertexAttribArray(pos)
	defer gl.DisableVertexAttribArray(pos)
	gl.VertexAttribPointer(pos, 2, gl.Float32, false, 0, 0)
	
	offset, _ := gl.GetUniformLocation(prog, "offset")
	
	clock := time.Tick(time.Second / 60)
	start := time.Now()
	
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
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		dx, dy := computeOffset(start)
		gl.Uniformf(offset, dx, dy)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		win.Flip()
		win.CheckEvent()
	}
}

func computeOffset(start time.Time) (dx float32, dy float32) {
	π := float64(math.Pi)
	period := time.Second * 2
	scale := 2*π / period.Seconds()
	elapsed := time.Since(start)
	pos := (elapsed % period).Seconds()
	
	dx = float32(math.Cos(pos * scale) / 2)
	dy = float32(math.Sin(pos * scale) / 2)
	return
}
