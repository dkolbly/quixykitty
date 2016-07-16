package main

import (
	"encoding/binary"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
	"log"
)

type Quix struct {
	X0, Y0    float32
	X1, Y1    float32
	age       float32
	color     float32
	lineIndex int
}

type Game struct {
	quixen      []Quix
	quixLineBuf gl.Buffer
	quixShader  gl.Program
	position    gl.Attrib
	color       gl.Uniform
}

func NewGame() *Game {
	g := &Game{}
	g.SpawnQuixen()
	return g
}

func (g *Game) SpawnQuixen() {
	q := Quix{
		X0: 10,
		Y0: 10,
		X1: 200,
		Y1: 50,
	}
	g.quixen = append(g.quixen, q)
}

func (g *Game) start(glctx gl.Context) error {
	g.quixLineBuf = glctx.CreateBuffer()

	shader, err := glutil.CreateProgram(
		glctx,
		quixVertexShader,
		quixFragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return err
	}
	g.quixShader = shader

	g.position = glctx.GetAttribLocation(shader, "position")
	g.color = glctx.GetUniformLocation(shader, "color")

	glctx.BindBuffer(gl.ARRAY_BUFFER, g.quixLineBuf)

	lines := f32.Bytes(
		binary.LittleEndian,
		//10.0, 10.0, 0, 100.0, 20.0, 0,
		0.1, 0.2, 0,   0.5, 0.5, 0,
		13.0, 15.0, 0, 130.0, 40.0, 0,
		16.0, 20.0, 0, 160.0, 60.0, 0,
	)

	//glctx.BufferSubData(gl.ARRAY_BUFFER, 0, lines)
	glctx.BufferData(gl.ARRAY_BUFFER, lines, gl.STATIC_DRAW)
	return nil
}

func (g *Game) paint(glctx gl.Context, sz size.Event) {
	glctx.UseProgram(g.quixShader)

	glctx.Uniform4f(g.color, 1, 1, 1, 1)

	// configure the 'position' attribute to use the quixLineBuf
	glctx.BindBuffer(gl.ARRAY_BUFFER, g.quixLineBuf)
	glctx.EnableVertexAttribArray(g.position)
	glctx.VertexAttribPointer(g.position, 3, gl.FLOAT, false, 0, 0)

	glctx.DrawArrays(gl.LINES, 0, 3)
	glctx.DisableVertexAttribArray(g.position)
}

const quixVertexShader = `#version 100
attribute vec4 position;
void main() {
   gl_Position = position;
}`

const quixFragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
