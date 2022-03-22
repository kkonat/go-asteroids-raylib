package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ship struct {
	shape     *shape
	m         *motion
	thr       V2
	mass      float64
	fuel      float64
	col       rl.Color
	isSliding bool
}

const S = 16

var shipShape = []V2{{-S / 2, -S}, {0, S}, {S / 2, -S}}

func newShip(posX, posY, mass, fuel float64) *ship {

	s := new(ship)
	s.shape = newShape(shipShape)
	s.m = newMotion()
	s.m.pos.x, s.m.pos.y = posX, posY
	s.col = rl.White
	s.mass = mass
	s.fuel = fuel
	return s
}
func (s *ship) Draw() {

	s.shape.Draw(s.m, rl.DarkGray, s.col)

	// draw thruster
	rl.DrawLineEx(
		rl.Vector2{
			X: float32(s.m.pos.x - s.thr.x*200), Y: float32(s.m.pos.y - s.thr.y*200)},
		rl.Vector2{
			X: float32(s.m.pos.x - s.thr.x*400), Y: float32(s.m.pos.y - s.thr.y*400)},
		4.1, rl.Orange)

	s.m.speed = V2MulA(s.m.speed, 0.9975)
	s.m.rotSpeed *= 0.97

}
func (s *ship) thrust(fuelCons float64) {
	//if s.fuel > 0 {
	force := fuelCons
	//	s.fuel -= fuelCons
	a := force * 100 / (s.mass + s.fuel)
	s.thr = V2MulA(cs(s.m.rot), a)
	s.m.speed = V2Add(s.m.speed, s.thr)
}

func (s *ship) rotate(dSpeed float64) {
	s.m.rotSpeed += dSpeed
}
