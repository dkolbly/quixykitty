package main

import (
	"log"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/sprite"
	"image"
	_ "image/png"
)

func loadTextures(eng sprite.Engine) []sprite.SubTex {
	a, err := asset.Open("sprite.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	const n = 48
	// The +1's and -1's in the rectangles below are to prevent colors from
	// adjacent textures leaking into a given texture.
	// See: http://stackoverflow.com/questions/19611745/opengl-black-lines-in-between-tiles
	ret := []sprite.SubTex{}
	for i := 0; i < 5; i++ {
		s := sprite.SubTex{t, image.Rect(n*i+1, 0, n*(i+1)-1, n)}
		ret = append(ret, s)
	}
	return ret
}
