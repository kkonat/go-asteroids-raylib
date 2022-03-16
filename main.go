package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type missile struct {
	shape shape
}

func newMissile() *missile {
	m := new(missile)

	m.shape = *newShape()

	//m.shape.addxy(0, 0)
	//m.shape.addxy(0, 20)
	m.shape.addxy(-1.25, 12.)
	m.shape.addxy(-1.2, 3.12)
	m.shape.addxy(-4.4, 0)
	m.shape.addxy(4.4, 0)
	m.shape.addxy(1.25, 3.12)
	m.shape.addxy(1.25, 12.5)
	m.shape.addxy(0, 20)

	return m
}

func (m *missile) launch(ship *Ship) {
	m.shape.pos = ship.shape.pos
	spd := ship.shape.speed
	m.shape.rot = ship.shape.rot
	dir := cs(m.shape.rot)
	dir = V2add(dir, spd)
	m.shape.speed = V2V(dir, 2.0)
}

func (m *missile) Draw() {
	m.shape.Draw(rl.Brown, rl.Brown)
}

func main() {

	game := newGame(1440, 720)

	for !rl.WindowShouldClose() {
		if !game.sm.isPlaying(0) {
			game.sm.play(game.sm.sSpace)
			//			fmt.Println("start backg sound")
		}

		if rl.IsKeyPressed('Q') {
			game.sm.playM(game.sm.sOinx)
			game.ship.shape.speed = V2{0, 0}
			game.ship.shape.pos = V2{720, 360}
		}
		if rl.IsKeyPressed('M') {
			if !game.sm.mute {
				game.sm.stop(0)
			}
			game.sm.mute = !game.sm.mute

		}
		if rl.IsKeyDown('A') {
			game.ship.rotate(-.2)
		}
		if rl.IsKeyPressed('S') {
			game.sm.play(game.sm.sThrust)
			game.ship.Slides = false
		}
		if rl.IsKeyDown('S') {
			game.ship.thrust(1.0)
		}
		if rl.IsKeyReleased('S') {
			game.ship.thrust(0)
			game.ship.Slides = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyPressed('W') {
			game.sm.play(game.sm.sThrust)
			game.ship.Slides = false
		}
		if rl.IsKeyDown('W') {
			game.ship.thrust(2.0)
		}
		if rl.IsKeyReleased('W') {
			game.ship.thrust(0)
			game.ship.Slides = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyDown('D') {
			game.ship.rotate(.2)
		}
		if rl.IsKeyPressed(rl.KeyLeftControl) {
			m := newMissile()
			m.launch(game.ship)
			game.missiles = append(game.missiles, *m)

			fmt.Println("launch")
		}

		game.drawGame()

	}
	game.finalize()

}
