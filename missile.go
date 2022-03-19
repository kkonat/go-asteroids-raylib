package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var missileShape = []V2{{-1.25, 12.}, {-1.2, 3.12}, {-4.4, 0}, {4.4, 0}, {1.25, 3.12}, {1.25, 12.5}, {0, 20}}

//

type missile struct {
	shape *shape
	m     *motion
}

func launchMissile(game *game) {

	m := new(missile)
	m.m = newMotion()
	m.shape = newShape(missileShape)

	m.m.pos = game.ship.m.pos.Add(game.ship.m.speed)
	spd := V2len(game.ship.m.speed)
	m.m.rot = game.ship.m.rot
	m.m.rotM = newM22rot(game.ship.m.rot)

	dir := cs(m.m.rot)
	m.m.speed = dir.MulA(spd + 2.0)

	game.missiles[game.missilesNo] = m
	game.missilesNo++
}

func (m missile) Draw() {
	m.shape.Draw(m.m, rl.Brown, rl.Brown)
}
