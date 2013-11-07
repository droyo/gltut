package main

import (
	"log"
	"time"
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
uniform float period;
uniform float time;

void main() {
	float scale = 3.14159 * 2 / period;
	float cur = mod(time, period);
	vec2 offset = vec2(cos(cur * scale) / 2, sin(cur * scale) / 2);
	gl_Position = vec4(position + offset, 0, 1);
}
`)

var fragShader = []byte(
`#version 150

out vec4 outColor;

uniform float fragPeriod;
uniform float time;

const vec4 firstColor = vec4(1, 1, 1, 1);
const vec4 secondColor = vec4(0, 1, 0, 1);

void main() {
	float cur = mod(time, fragPeriod);
	float curLerp = cur / fragPeriod;
	outColor = mix(firstColor, secondColor, curLerp);
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
	
	glTime, _ := gl.GetUniformLocation(prog, "time")
	glPeriod, _ := gl.GetUniformLocation(prog, "period")
	glFragPeriod, _ := gl.GetUniformLocation(prog, "fragPeriod")
	
	clock := time.Tick(time.Second / 60)
	gl.Uniformf(glPeriod, float32((time.Second * 4).Seconds()))
	gl.Uniformf(glFragPeriod, float32((time.Second * 2).Seconds()))
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
		gl.Uniformf(glTime, float32(time.Since(start).Seconds()))
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		win.Flip()
		win.CheckEvent()
	}
}
