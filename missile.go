package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var missileShape = []V2{{-1.25, 12.}, {-1.2, 3.12}, {-4.4, 0}, {4.4, 0}, {1.25, 3.12}, {1.25, 12.5}, {0, 20}}

//

type missile struct {
	*shape
	motion
}

func launchMissile(game *game) {

	m := new(missile)

	m.shape = newShape(missileShape)

	m.pos = game.ship.pos.Add(game.ship.speed)
	spd := V2len(game.ship.speed)
	m.rot = game.ship.rot

	dir := cs(m.rot)
	m.speed = dir.MulA(spd + 2.0)

	game.missiles = append(game.missiles, m)
}

func (m missile) Draw() {
	m.shape.Draw(m.pos, m.rot, rl.Brown, rl.Brown)
}
