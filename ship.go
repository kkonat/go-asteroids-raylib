package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ship struct {
	shape  shape
	thr    V2
	mass   float64
	fuel   float64
	col    rl.Color
	Slides bool
}

const S = 16

var shipShape = []V2{{-S / 2, -S}, {0, S}, {S / 2, -S}}

func newShip(posX, posY, mass, fuel float64) *Ship {

	s := new(Ship)
	s.shape = *newShape(shipShape)
	//s.shape.addPoints(shipShape)

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
			X: float32(s.shape.pos.x - s.thr.x*200), Y: float32(s.shape.pos.y - s.thr.y*200)},
		rl.Vector2{
			X: float32(s.shape.pos.x - s.thr.x*400), Y: float32(s.shape.pos.y - s.thr.y*400)},
		4.1, rl.Orange)

	s.shape.speed = V2MulA(s.shape.speed, 0.9975)

}
func (s *Ship) thrust(fuelCons float64) {
	//if s.fuel > 0 {
	force := fuelCons
	//	s.fuel -= fuelCons
	a := force * 100 / (s.mass + s.fuel)
	s.thr = V2MulA(cs(s.shape.rot), a)
	s.shape.speed = V2Add(s.shape.speed, s.thr)
}

func (s *Ship) rotate(dSpeed float64) {
	s.shape.rotSpeed += dSpeed

}
