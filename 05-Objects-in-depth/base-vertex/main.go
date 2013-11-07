// base-vertex renders an object using DrawElementsBaseVertex.
// It is an implementation of http://arcsynthesis.org/gltut/Positioning/Tut05%20Optimization%20Base%20Vertex.html
package main

import (
	"log"
	"aqwari.net/exp/gl"
	"aqwari.net/exp/display"
)

var config = display.Config{
	"Title":          "Overlap No Depth",
	"Geometry":       "500x500",
	"OpenGL Version": "3.2",
}

var vertShader = []byte(
`#version 150

in vec4 position;
in vec4 color;

smooth out vec4 theColor;

uniform vec3 offset;
uniform mat4 perspectiveMatrix;

void main()
{
	vec4 camera = position + vec4(offset, 0);
	gl_Position = perspectiveMatrix * camera;
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
		Grey  = []float32{0.80, 0.80, 0.80, 0.80}
		Brown = []float32{0.50, 0.50, 0.00, 1.00}
	)
	
	vertexData := []float32{
		LeftExtent, TopExtent, RearExtent,
		LeftExtent, MiddleExtent, FrontExtent,
		RightExtent, MiddleExtent, FrontExtent,
		RightExtent, MiddleExtent, RearExtent,

		LeftExtent, BottomExtent, RearExtent,
		LeftExtent, MiddleExtent, FrontExtent,
		RightExtent, MiddleExtent, FrontExtent,
		RightExtent, BottomExtent, RearExtent,

		LeftExtent, TopExtent, RearExtent,
		LeftExtent, MiddleExtent, FrontExtent,
		LeftExtent, BottomExtent, RearExtent,

		RightExtent, TopExtent, RearExtent,
		RightExtent, MiddleExtent, FrontExtent,
		RightExtent, BottomExtent, RearExtent,

		LeftExtent, BottomExtent, RearExtent,
		LeftExtent, TopExtent, RearExtent,
		RightExtent, TopExtent, RearExtent,
		RightExtent, BottomExtent, RearExtent,

		//Object 2 positions
		TopExtent, RightExtent, RearExtent,
		MiddleExtent, RightExtent, FrontExtent,
		MiddleExtent, LeftExtent, FrontExtent,
		TopExtent, LeftExtent, RearExtent,

		BottomExtent, RightExtent, RearExtent,
		MiddleExtent, RightExtent, FrontExtent,
		MiddleExtent, LeftExtent, FrontExtent,
		BottomExtent, LeftExtent, RearExtent,

		TopExtent, RightExtent, RearExtent,
		MiddleExtent, RightExtent, FrontExtent,
		BottomExtent, RightExtent, RearExtent,

		TopExtent, LeftExtent, RearExtent,
		MiddleExtent, LeftExtent, FrontExtent,
		BottomExtent, LeftExtent, RearExtent,

		BottomExtent, RightExtent, RearExtent,
		TopExtent, RightExtent, RearExtent,
		TopExtent, LeftExtent, RearExtent,
		BottomExtent, LeftExtent, RearExtent,
	}
	
	// Object 1 colors
	for _, col := range [...][]float32{Green, Blue, Red, Grey, Brown} {
		for i := 0 ; i < 4 ; i++ {
			vertexData = append(vertexData, col...)
		}
	}
	
	// Object 2 colors
	for _, col := range [...][]float32{Red, Brown, Blue, Green, Grey} {
		for i := 0 ; i < 4 ; i++ {
			vertexData = append(vertexData, col...)
		}
	}
	
	indices := []uint16 {
		0, 2, 1,
		3, 2, 0,

		4, 5, 6,
		6, 7, 4,

		8, 9, 10,
		11, 13, 12,

		14, 16, 15,
		17, 16, 14,
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
	err = gl.BufferData(gl.ARRAY_BUFFER, vertexData, gl.STATIC_DRAW)
	if err != nil {
		log.Fatal(err)
	}
	
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf[1])
	err = gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)
	if err != nil {
		log.Fatal(err)
	}
	
	pos, _ := gl.GetAttribLocation(prog, "position")
	col, _ := gl.GetAttribLocation(prog, "color")
	vao := gl.GenVertexArrays(1)
	
	// Object 1
	gl.BindVertexArray(vao[0])
	gl.EnableVertexAttribArray(pos)
	gl.EnableVertexAttribArray(col)
	gl.VertexAttribPointer(pos, 3, gl.Float32, false, 0, 0)
	gl.VertexAttribPointer(col, 4, gl.Float32, false, 0, 4 * 3 * 36)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf[1])
	
	offset, _ := gl.GetUniformLocation(prog, "offset")
	perspective, _ := gl.GetUniformLocation(prog, "perspectiveMatrix")
	const (
		frustum float32 = 1.0
		zNear float32 = 1.0
		zFar float32 = 3.0
	)
	matrix := [16]float32{
		0:  frustum,
		5:  frustum,
		10: (zFar + zNear) / (zNear - zFar),
		14: (2 * zFar * zNear) / (zNear - zFar),
		11: -1.0,
	}
	gl.UniformMatrix4fv(perspective, false, matrix[:])
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
Loop:
	for {
		select {
		case ev := <-win.Event:
			switch ev := ev.(type) {
			case display.KeyPress:
				if ev.Code == display.KeyEscape {
					break Loop
				}
			case display.Resize:
				matrix[0] = frustum / (float32(ev.Width) / float32(ev.Height))
				matrix[5] = frustum
				gl.Viewport(0, 0, ev.Width, ev.Height)
				gl.UniformMatrix4fv(perspective, false, matrix[:])
			}
		default:
			win.WaitEvent()
			continue
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		
		gl.Uniformf(offset, 0, 0, 0)
		gl.DrawElements(gl.TRIANGLES, len(indices), gl.Uint16, 0)
		
		gl.Uniformf(offset, 0, 0, -1)
		gl.DrawElementsBaseVertex(gl.TRIANGLES, len(indices),
			gl.Uint16, 0, 36/2)
		
		win.Flip()
	}
}
