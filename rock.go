package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const noRocks = 16

type Rock struct {
	shape  *shape
	radius float64
	mass   float64
}

func newRock(g *game) *Rock {

	r := new(Rock)

	r.randomize()
	r.shape.pos = V2{rand.Float64() * float64(g.sW), rand.Float64() * float64(g.sH)}
	r.shape.speed = V2{rand.Float64()*rSpeedMax*2.0 - rSpeedMax, rand.Float64()*rSpeedMax*2.0 - rSpeedMax}
	r.shape.rotSpeed = rand.Float64()*0.2 - 0.1

	return r
}
func (r *Rock) randomize() {
	r.radius = 10 + rand.Float64()*100
	n := 6 + rand.Intn(10) + int(r.radius/5)
	data := make([]V2, n)

	var step = 360 / float64(n)

	angle := 0.0
	for i := 0; i < n; i++ {
		angle += step + rand.Float64()*step/2 - step/4
		r1 := r.radius + rand.Float64()*r.radius/4 - r.radius/8
		p := cs(angle)
		data[i] = p.MulA(r1)
	}

	r.shape = newShape(data)

	r.shape.rotSpeed = rand.Float64()*1.5 - 0.75

}

func (r *Rock) Draw() {
	//rl.DrawCircle(int32(r.shape.pos.x), int32(r.shape.pos.y), r.radius, rl.ColorAlpha(rl.DarkGray, 0.2))
	//	rl.DrawLine(720,360,int32(r.shape.pos.x),int32(r.shape.pos.y),rl.DarkBlue)
	r.shape.Draw(rl.Black, rl.DarkGray)
}
