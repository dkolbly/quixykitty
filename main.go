// +build darwin linux windows

// An app that allows a kitty to capture territory while being harassed
// by lines
//
package main

import (
	"time"
	
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	//"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	//fps      *debug.FPS
)

var eng sprite.Engine
var scene *sprite.Node
var texs []sprite.SubTex
var startTime time.Time = time.Now()

func main() {
	app.Main(func(a app.App) {

		var glctx gl.Context
		var currentSize size.Event
		
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				currentSize = e
				// TODO: Screen reorientation?
				//touchX = float32(sz.WidthPx / 2)
				//touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting
					// as fast as we can (usually
					// 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, currentSize)
				a.Publish()
				// Drive the animation by preparing to
				// paint the next frame after this one
				// is shown.
				a.Send(paint.Event{})
			case touch.Event:
				gg.touch.X = int(e.X)
				gg.touch.Y = int(e.Y)
			}
		}
	})
}

var gg *Game

func onStart(glctx gl.Context) {
	gg = NewGame()
	gg.start(glctx)

	images = glutil.NewImages(glctx)
	//fps = debug.NewFPS(images)
	eng = glsprite.Engine(images)

	// exprimental sprite stuff
	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	// add a node
	thingy := &sprite.Node{
		Arranger: gg.kitty,
	}
	eng.Register(thingy)
	scene.AppendChild(thingy)
	texs = loadTextures(eng)
}

func onStop(glctx gl.Context) {
	gg.stop(glctx)
	//fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(0.7, 0.8, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)

	gg.paint(glctx, sz, now)

	//fps.Draw(sz)
}
