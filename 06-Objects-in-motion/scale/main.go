package main

import (
	"log"
	"time"
	"math"
	"aqwari.net/exp/gl"
	"aqwari.net/exp/display"
)

var config = display.Config{
	"Title":          "Scaling",
	"Geometry":       "500x500",
	"OpenGL Version": "3.2",
}

var vertShader = []byte(
`#version 150

in vec4 position;
in vec4 color;

smooth out vec4 theColor;

uniform mat4 cameraToClipMatrix;
uniform mat4 modelToCameraMatrix;

void main()
{
	vec4 cameraPos = modelToCameraMatrix * position;
	gl_Position = cameraToClipMatrix * cameraPos;
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

func CalcFrustum(deg float64) float32 {
	π := float64(math.Pi)
	deg2rad := π / 2 / 360
	return float32(1 / math.Tan(deg * deg2rad / 2))
}

func UpdateOval(elapsed time.Duration, mat []float32) {
	π := float64(math.Pi)
	period := time.Second * 3
	scale := π * 2 / period.Seconds()
	
	pos := (elapsed % period).Seconds()
	mat[3] = float32(math.Cos(pos * scale) * 4)
	mat[7] = float32(math.Sin(pos * scale) * 6)
	mat[11] = -20
}

func CalcLerpFactor(elapsed, period time.Duration) float32 {
	if pos := elapsed % period; pos < time.Second / 2 {
		return float32((time.Second - pos).Seconds())
	} else {
		return float32((pos * 2 * time.Second).Seconds())
	}
}

func UpdateUniformScale(elapsed time.Duration, mat []float32) {
	
	pos := elapsed % period
	if pos < time.Second / 2 {
		time.Second - pos
	}
}

func UpdateCircle(elapsed time.Duration, mat []float32) {
	π := float64(math.Pi)
	period := time.Second * 12
	scale := π * 2 / period.Seconds()
	
	pos := (elapsed % period).Seconds()
	mat[3] = float32(math.Cos(pos * scale) * 5)
	mat[7] = -3.5
	mat[11] = float32(math.Sin(pos * scale) * 5 - 20)
}

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
	gl.ClearDepth(1)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthMask(true)
	gl.DepthRange(0, 1)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	
	const (
		RightExtent = 0.8
		LeftExtent = -RightExtent
		TopExtent = 0.2
		MiddleExtent = 0.0
		BottomExtent = -TopExtent
		FrontExtent = -1.25
		RearExtent = -1.75		
	)
	
	var (
		Green = []float32{0.75, 0.75, 1.00, 1.00}
		Blue  = []float32{0.00, 0.50, 0.00, 1.00}
		Red   = []float32{1.00, 0.00, 0.00, 1.00}
		Brown = []float32{0.50, 0.50, 0.00, 1.00}
	)
	
	vertexData := []float32{
		+1, +1, +1,
		-1, -1, +1,
		-1, +1, -1,
		+1, -1, -1,
		
		-1, -1, -1,
		+1, +1, -1,
		+1, -1, +1,
		-1, +1, +1,
	}
	
	// Object 1 colors
	for _, col := range [...][]float32{Green, Blue, Red, Brown} {
		vertexData = append(vertexData, col...)
	}
	
	// Object 2 colors
	for _, col := range [...][]float32{Green, Blue, Red, Brown} {
		vertexData = append(vertexData, col...)
	}
	
	indices := []uint16 {
		0, 1, 2,
		1, 0, 3,
		2, 3, 0,
		3, 2, 1,

		5, 4, 6,
		4, 5, 7,
		7, 6, 4,
		6, 7, 5,
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
	
	buf := gl.GenBuffers(2)
	defer gl.DeleteBuffers(buf)
	
	gl.BindBuffer(gl.ARRAY_BUFFER, buf[0])
	gl.BufferData(gl.ARRAY_BUFFER, vertexData, gl.STATIC_DRAW)
	
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)
	
	pos, _ := gl.GetAttribLocation(prog, "position")
	col, _ := gl.GetAttribLocation(prog, "color")
	vao := gl.GenVertexArrays(1)
	
	// Object 1
	gl.BindVertexArray(vao[0])
	gl.EnableVertexAttribArray(pos)
	gl.EnableVertexAttribArray(col)
	gl.VertexAttribPointer(pos, 3, gl.Float32, false, 0, 0)
	gl.VertexAttribPointer(col, 4, gl.Float32, false, 0, 4 * 3 * 8)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf[1])
	
	offset, _ := gl.GetUniformLocation(prog, "modelToCameraMatrix")
	perspective, _ := gl.GetUniformLocation(prog, "cameraToClipMatrix")
	
	var (
		frustum float32 = CalcFrustum(125)
		zNear float32 = 1
		zFar float32 = 61
	)
	perspectiveMatrix := [16]float32{
		0:  frustum,
		5:  frustum,
		10: (zFar + zNear) / (zNear - zFar),
		14: (2 * zFar * zNear) / (zNear - zFar),
		11: -1.0,
	}
	identity := []float32 {
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	var (
		stationary = identity
		circular = make([]float32, 16)
		ovular = make([]float32, 16)
	)
	copy(circular, identity)
	copy(ovular, identity)
	stationary[3] = 0
	stationary[7] = 0
	stationary[11] = -20

	gl.UniformMatrix4fv(perspective, false, perspectiveMatrix[:])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	
	clock := time.Tick(time.Second / 60)
	start := time.Now()
Loop:
	for _ = range clock {
EventRead:
		for {
			select {
			case ev := <-win.Event:
				switch ev := ev.(type) {
				case display.KeyPress:
					if ev.Code == display.KeyEscape {
						break Loop
					}
				case display.Resize:
					perspectiveMatrix[0] = frustum / (float32(ev.Width) / float32(ev.Height))
					perspectiveMatrix[5] = frustum
					gl.Viewport(0, 0, ev.Width, ev.Height)
					gl.UniformMatrix4fv(perspective, false, perspectiveMatrix[:])
				}
			default:
				win.CheckEvent()
				break EventRead
			}
		}
		elapsed := time.Since(start)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		
		gl.UniformMatrix4fv(offset, true, stationary)
		gl.DrawElements(gl.TRIANGLES, len(indices), gl.Uint16, 0)
		
		UpdateCircle(elapsed, circular)
		gl.UniformMatrix4fv(offset, true, circular)
		gl.DrawElements(gl.TRIANGLES, len(indices), gl.Uint16, 0)
		
		UpdateOval(elapsed, ovular)
		gl.UniformMatrix4fv(offset, true, ovular)
		gl.DrawElements(gl.TRIANGLES, len(indices), gl.Uint16, 0)
		win.Flip()
	}
}
