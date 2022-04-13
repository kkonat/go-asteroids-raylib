package main

import (
	v "rlbb/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ship struct {
	shape *shape
	motion
	thr        V2
	mass       float64
	energy     float64
	shields    float64
	cash       int
	col        rl.Color
	isSliding  bool
	cycle      uint8
	forceField bool
	destroyed  bool
	light      *OmniLight
}

const S = 16

var shipShape = []V2{{X: S / 2, Y: -S}, {X: 0, Y: S}, {X: -S / 2, Y: -S}}

func newShip(posX, posY, mass, fuel float64) *ship {

	s := new(ship)
	s.shields = 100
	s.energy = 1000

	s.destroyed = false
	s.shape = newShape(shipShape)

	s.pos.X, s.pos.Y = posX, posY
	// s.m.rot = rnd() * 360
	// s.m.speed = cs(s.m.rot)

	s.col = rl.White
	s.mass = mass
	s.energy = fuel

	s.light = &OmniLight{s.pos, Color{0, 0, 0, 1}, 10}

	return s
}

func (s *ship) chargeUp() {
	if !s.destroyed {
		dist := V2{X: 1655, Y: 400}.Sub(s.pos).Len()
		chUp := 16 / dist
		s.energy += chUp
		if s.energy > 1000 {
			s.energy = 1000
		}
		if s.shields < 100 {
			s.shields += 0.001
		}
	}
}

func (s *ship) Move(dt float64) {
	s.motion.Move(dt)
	s.light.Pos = s.pos.Sub(s.thr.Norm().MulA(24))
	s.speed = s.speed.MulA(0.9975)
	s.rotSpeed *= 0.97
}
func (s *ship) Draw() {
	if !s.destroyed {
		// draw ship
		s.shape.Draw(s.pos, s.rot, rl.Black, s.col)

		thr := s.thr.Len()
		// draw flame
		disturb := _noise2D(s.cycle * 4).MulA(6).SubA(3)
		p1 := s.pos.Sub(s.thr.Norm().MulA(16))
		p2 := p1.Sub(s.thr.MulA(200)).Add(disturb)

		n := _noise1D(s.cycle)
		c := _colorBlendA(n, rl.Yellow, rl.Red)

		s.light.Col = newColorRGBint(c.R, c.G, c.B)
		s.light.Strength = thr * 500

		_lineThick(p1, p2, 4.1, c)
		// animate

		// animate noise
		s.cycle++
	}
}
func (s *ship) thrust(fuelCons float64) {
	if s.energy > fuelCons {
		s.energy -= fuelCons
		force := fuelCons

		a := force * 0.1
		s.thr = v.Cs(s.rot).MulA(a)
		s.speed = s.speed.Add(s.thr)
	}
}

func (s *ship) rotate(dSpeed float64) {
	s.rotSpeed += dSpeed
}
