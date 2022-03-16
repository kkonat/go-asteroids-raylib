package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ship struct {
	shape  shape
	thr    V2
	mass   float32
	fuel   float32
	col    rl.Color
	Slides bool
}

func newShip(posX, posY, mass, fuel float32) *Ship {
	const S = 16
	s := new(Ship)
	s.shape = *newShape()
	s.shape.add(V2{-S / 2, -S})
	s.shape.add(V2{0, S})
	s.shape.add(V2{S / 2, -S})
	s.shape.pos.x, s.shape.pos.x = posX, posY
	s.col = rl.White
	s.mass = mass
	s.fuel = fuel
	return s
}
func (s *Ship) Draw() {

	s.shape.Draw(rl.DarkGray, s.col)

	// draw thruster
	rl.DrawLineEx(
		rl.Vector2{
			X: s.shape.pos.x - s.thr.x*200, Y: s.shape.pos.y - s.thr.y*200},
		rl.Vector2{
			X: s.shape.pos.x - s.thr.x*400, Y: s.shape.pos.y - s.thr.y*400},
		4.1, rl.Orange)

	s.shape.speed = V2V(s.shape.speed, 0.9975)

}
func (s *Ship) thrust(fuelCons float32) {
	//if s.fuel > 0 {
	force := fuelCons
	//	s.fuel -= fuelCons
	a := force * 100 / (s.mass + s.fuel)
	s.thr = V2V(cs(s.shape.rot), a)
	s.shape.speed = V2add(s.shape.speed, s.thr)
}

func (s *Ship) rotate(dSpeed float32) {
	s.shape.rotSpeed += dSpeed

}
