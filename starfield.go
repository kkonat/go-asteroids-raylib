package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type star struct {
	y        int32
	x, speed float32
	r        uint8
}

type starfield struct {
	stars []star
	w, h  int32
}

const starsNo = 1000

func newStarfield(w, h int32) *starfield {
	sf := new(starfield)
	sf.w, sf.h = w, h
	var s star
	sf.stars = make([]star, starsNo)
	for i := 0; i < starsNo; i++ {
		s.x = rand.Float32() * float32(w)
		s.y = rand.Int31n(h)
		s.speed = rand.Float32() * 5.0
		s.r = uint8(rand.Int31n(40) * 3)
		sf.stars[i] = s
	}

	return sf
}
func (sf *starfield) draw() {
	var x float32
	c := rl.NewColor(0, 0, 0, 255)
	rl.DrawCircleGradient(1655, 400, 500, rl.NewColor(60, 40, 210, 255), rl.Black)
	rl.DrawCircleLines(1655, 400, 400, rl.NewColor(40, 10, 140, 255))
	rl.DrawCircle(1655, 400, 400, rl.NewColor(20, 0, 70, 255))

	for i, s := range sf.stars {
		c.B = uint8(s.speed * 25)
		c.R = s.r

		//rl.DrawPixel(int32(s.x), int32(s.y), c)
		rl.DrawCircle(int32(s.x), int32(s.y), 1+float32(s.speed/3.0), c)
		x = s.x + s.speed*0.025
		if x > float32(sf.w) {
			x -= float32(sf.w)
		}
		sf.stars[i].x = x
	}
}
