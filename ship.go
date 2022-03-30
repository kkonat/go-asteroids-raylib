package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ship struct {
	shape     *shape
	m         *motion
	thr       V2
	mass      float64
	energy    float64
	shields   float64
	missiles  int
	cash      int
	col       rl.Color
	isSliding bool
	cycle     uint8
	destroyed bool
}

const S = 16

var shipShape = []V2{{-S / 2, -S}, {0, S}, {S / 2, -S}}

func newShip(posX, posY, mass, fuel float64) *ship {

	s := new(ship)
	s.shields = 100
	s.missiles = 100
	s.energy = 1000

	s.destroyed = false
	s.shape = newShape(shipShape)
	s.m = newMotion()
	s.m.pos.x, s.m.pos.y = posX, posY
	s.m.rot = rnd() * 360
	s.m.speed = cs(s.m.rot)

	s.col = rl.White
	s.mass = mass
	s.energy = fuel

	return s
}

func (s *ship) chargeUp() {
	if !s.destroyed {
		dist := V2{1655, 400}.Sub(s.m.pos).Len()
		chUp := 16 / dist
		if s.energy < 1000-chUp {
			s.energy += chUp
			if s.energy >= 1000 {
				s.energy = 1000
			}
		}
		if s.shields < 100 {
			s.shields += 0.001
		}
	}
}

func (s *ship) Draw() {
	if !s.destroyed {
		// draw ship
		s.shape.Draw(s.m, rl.DarkGray, s.col)

		// draw flame
		disturb := _noise2D(s.cycle * 4).MulA(6).SubA(3)
		p1 := s.m.pos.Sub(s.thr.Norm().MulA(16))
		p2 := p1.Sub(s.thr.MulA(200)).Add(disturb)

		n := _noise1D(s.cycle)
		c := _colorBlendA(n, rl.Yellow, rl.Red)
		_lineThick(p1, p2, 4.1, c)

		// animate
		s.m.speed = V2MulA(s.m.speed, 0.9975)
		s.m.rotSpeed *= 0.97

		// animate noise
		s.cycle++
	}
}
func (s *ship) thrust(fuelCons float64) {
	if s.energy > fuelCons {
		s.energy -= fuelCons
		force := fuelCons

		a := force * 0.1
		s.thr = cs(s.m.rot).MulA(a)
		s.m.speed = s.m.speed.Add(s.thr)
	}
}

func (s *ship) rotate(dSpeed float64) {
	s.m.rotSpeed += dSpeed
}
