package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rock struct {
	shape  shape
	radius float32
	mass   float32
}

func newRock(pos, speed V2, rotSpeed float32) *Rock {

	r := new(Rock)
	r.Randomize()

	r.shape.pos = pos
	r.shape.speed = speed
	//	r.shape.rotSpeed = rotSpeed

	return r
}
func (r *Rock) Randomize() {
	r.shape = *newShape()

	r.shape.rotSpeed = rand.Float32()*1.5 - 0.75
	r.radius = 10 + rand.Float32()*100
	n := 6 + rand.Intn(10) + int(r.radius/5)

	var step float32 = 360 / float32(n)
	angle := float32(0)
	for i := 0; i < n; i++ {
		angle += step + rand.Float32()*step/2 - step/4
		r1 := r.radius + rand.Float32()*r.radius/4 - r.radius/8
		p := cs(angle)
		v := V2V(p, r1)
		r.shape.add(v)
	}
}

func (r *Rock) Draw() {
	//rl.DrawCircle(int32(r.shape.pos.x), int32(r.shape.pos.y), r.radius, rl.ColorAlpha(rl.DarkGray, 0.2))

	r.shape.Draw(rl.Black, rl.DarkGray)
}
