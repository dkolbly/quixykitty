// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux windows

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	//program  gl.Program
	//position gl.Attrib
	//offset   gl.Uniform
	//color    gl.Uniform
	//buf      gl.Buffer

	green  float32
	touchX float32
	touchY float32
)

var eng sprite.Engine
var scene *sprite.Node
var texs []sprite.SubTex
var startTime time.Time = time.Now()

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
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
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting
					// as fast as we can (usually
					// 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to
				// paint the next frame after this one
				// is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchX = e.X
				touchY = e.Y
				gg.touch.X = e.X
				gg.touch.Y = e.Y
			}
		}
	})
}

var gg *Game

func onStart(glctx gl.Context) {
	gg = NewGame()
	gg.start(glctx)

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
	green = 0.5
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
	//glctx.DeleteProgram(program)
	//glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(0.7, 0.8, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	eng.Render(scene, now, sz)

	gg.paint(glctx, sz, now)

	fps.Draw(sz)
}
