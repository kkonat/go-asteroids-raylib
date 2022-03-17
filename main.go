package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
		if rl.IsKeyDown('A') { // rotate left
			game.ship.rotate(-.2)
		}
		if rl.IsKeyPressed('S') { // small thrust
			game.sm.play(game.sm.sThrust)
			game.ship.Slides = false
		}
		if rl.IsKeyDown('S') {
			game.ship.thrust(1.0)
		}
		if rl.IsKeyReleased('S') { // -----
			game.ship.thrust(0)
			game.ship.Slides = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyPressed('W') { // big thrust
			game.sm.play(game.sm.sThrust)
			game.ship.Slides = false
		}
		if rl.IsKeyDown('W') {
			game.ship.thrust(2.0)
		}
		if rl.IsKeyReleased('W') { // -----
			game.ship.thrust(0)
			game.ship.Slides = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyDown('D') { // rotate right
			game.ship.rotate(.2)
		}
		if rl.IsKeyPressed(rl.KeyLeftControl) { // fire
			if game.missilesNo < len(game.missiles) {
				launchMissile(game)
			}

		}
		if rl.IsKeyDown(rl.KeyCapsLock) { // slow down rotation
			game.ship.shape.rotSpeed *= 0.9
		}

		game.drawGame()

	}
	game.finalize()

}
