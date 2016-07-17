package main

import (
	"image"
	"encoding/binary"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/gl"
	"log"
	"math"
	"math/rand"

	"github.com/dkolbly/quixykitty/polygon"
)

type QuixEnd struct {
	X, Y   float32
	Vx, Vy float32
}

type Quix struct {
	A, B        QuixEnd
	age         float32
	color       float32
	offset      int
	head        int
	reservation int
	length      int
	nextBump    clock.Time
}

func (e *QuixEnd) bounce() {
	if e.X+e.Vx > 100 {
		e.Vx = -(rand.Float32() * 10) + 0.3
	} else if e.X+e.Vx < -100 {
		e.Vx = (rand.Float32() * 10) + 0.3
	}
	if e.Y+e.Vy > 100 {
		e.Vy = -(rand.Float32() * 10) + 0.3
	} else if e.Y+e.Vy < -100 {
		e.Vy = (rand.Float32() * 10) + 0.3
	}
}

func (e *QuixEnd) bump() {
	e.X += e.Vx
	e.Y += e.Vy

	// TODO: bounce
}

const bumpRate = 4

func (q *Quix) bump(glctx gl.Context, g *Game, t clock.Time) {
	if t < q.nextBump {
		return
	}
	//log.Printf("bump t=%d", t)
	// add a new segment
	if q.head >= q.reservation {
		q.head = 0
	}
	(&q.A).bounce()
	(&q.A).bump()
	(&q.B).bounce()
	(&q.B).bump()

	g.SetSegment(
		glctx,
		q.offset+q.head,
		q.A.X, q.A.Y,
		q.B.X, q.B.Y,
	)
	q.head++
	q.length++
	if q.length > q.reservation {
		q.length = q.reservation
	}

	if t > q.nextBump+bumpRate {
		q.nextBump = t + bumpRate
	} else {
		q.nextBump = q.nextBump + bumpRate
	}
}

type CapturedRegion struct {
	offset int
	length int
}

type Game struct {
	quixen      []Quix
	captured []CapturedRegion
	quixLineBuf gl.Buffer
	quixShader  gl.Program
	position    gl.Attrib
	color       gl.Uniform
	kitty       *KittySprite
	touch	    image.Point
}

func NewGame() *Game {
	g := &Game{}
	g.kitty = &KittySprite{P: QuixEnd{X: 100, Y: 200, Vx: 1, Vy: -1}}
	g.SpawnQuixen()

	// test polygon capture
	v := []image.Point{
		{20, 50},
		{20, 22},
		{48, 15},
		{70, 22},
		{65, 50},
		{52, 56},
		{45, 50},
		{40, 57},
		
	}
	l := polygon.Triangulate(v)
	log.Printf("Indexes: %#v", l)
	
	return g
}

func (g *Game) SetSegment(glctx gl.Context, i int, x0, y0, x1, y1 float32) {
	seg := f32.Bytes(binary.LittleEndian,
		x0, y0, 0,
		x1, y1, 0,
	)
	glctx.BindBuffer(gl.ARRAY_BUFFER, g.quixLineBuf)
	glctx.BufferSubData(gl.ARRAY_BUFFER, 4*6*i, seg)
}

func (g *Game) SpawnQuixen() {
	q := Quix{
		offset:      0,
		head:        0,
		length:      0,
		reservation: 20,
		A: QuixEnd{
			X:  0,
			Y:  0,
			Vx: 2,
			Vy: 3,
		},
		B: QuixEnd{
			X:  40,
			Y:  0,
			Vx: 3,
			Vy: 2,
		},
	}
	g.quixen = append(g.quixen, q)
}

func (g *Game) stop(glctx gl.Context) {
	glctx.DeleteProgram(g.quixShader)
	glctx.DeleteBuffer(g.quixLineBuf)
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

	lines := make([]byte, 4*6*20)
	/*lines := f32.Bytes(
		binary.LittleEndian,
		10.0, 20.0, 0, 50.0, 45.0, 0,
		13.0, 15.0, 0, 130.0, 40.0, 0,
		16.0, 20.0, 0, 160.0, 60.0, 0,
	)*/

	glctx.BufferData(gl.ARRAY_BUFFER, lines, gl.DYNAMIC_DRAW)
	//glctx.BufferSubData(gl.ARRAY_BUFFER, 0, lines)
	return nil
}

func (g *Game) paint(glctx gl.Context, sz size.Event, t clock.Time) {

	(&g.quixen[0]).bump(glctx, g, t)
	g.kitty.bump(g, t)
	
	glctx.UseProgram(g.quixShader)

	glctx.Uniform4f(g.color, 1, 1, 1, 1)

	// configure the 'position' attribute to use the quixLineBuf
	glctx.BindBuffer(gl.ARRAY_BUFFER, g.quixLineBuf)
	glctx.EnableVertexAttribArray(g.position)
	glctx.VertexAttribPointer(g.position, 3, gl.FLOAT, false, 0, 0)

	glctx.DrawArrays(gl.LINES, 0, g.quixen[0].length*3)
	glctx.DisableVertexAttribArray(g.position)
}

type KittySprite struct {
	P        QuixEnd
	nextBump clock.Time
}

func (ks *KittySprite) bump(g *Game, t clock.Time) {
	if t < ks.nextBump {
		return
	}
	ks.P.Vx = float32(g.touch.X) - ks.P.X
	ks.P.Vy = float32(g.touch.Y) - ks.P.Y
	veclen := math.Sqrt(float64(ks.P.Vx) * float64(ks.P.Vx) +
		float64(ks.P.Vy) * float64(ks.P.Vy))
	if veclen <= 0 {
		ks.P.Vx = 0
		ks.P.Vy = 0
	} else {
		ks.P.Vx /= float32(veclen)
		ks.P.Vy /= float32(veclen)
		if veclen > 100 {
			// high speed mode
			ks.P.Vx *= 2
			ks.P.Vy *= 2
		}
	}
		
	
	(&ks.P).bump()
	ks.nextBump = t+1
}

func (ks *KittySprite) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	eng.SetSubTex(n, texs[5])
	a := float64(t%120) / 120 * 2 * math.Pi

	sin := float32(48 * math.Sin(a))
	cos := float32(48 * math.Cos(a))

	// translate, scale, and rotate
	tsr := f32.Affine{{cos, sin, ks.P.X}, {-sin, cos, ks.P.Y}}
	var m f32.Affine
	(&m).Translate(&tsr, -0.5, -0.5)
	eng.SetTransform(n, m)
}

const quixVertexShader = `#version 100
attribute vec4 position;
void main() {
   gl_Position = vec4(position.x * 0.01, position.y * 0.01, 0, 1);
}`

const quixFragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
