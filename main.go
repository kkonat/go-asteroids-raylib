package main

import (
	"math/rand"
	"runtime"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	runtime.GOMAXPROCS(4)
	rand.Seed(time.Now().UnixNano())
	game := newGame(1440, 720)
	rl.DisableCursor()
	cursorEnabled := false

	for !rl.WindowShouldClose() {
		if !game.sm.isPlaying(0) {
			game.sm.play(game.sm.sSpace)
			//			fmt.Println("start backg sound")
		}

		if rl.IsKeyPressed('Q') {
			game.sm.playM(game.sm.sOinx)
			game.ship.m.speed = V2{0, 0}
			game.ship.m.pos = V2{720, 360}
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
			game.ship.isSliding = false
		}
		if rl.IsKeyDown('S') {
			game.ship.thrust(1.0)
		}
		if rl.IsKeyReleased('S') { // -----
			game.ship.thrust(0)
			game.ship.isSliding = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyPressed('W') { // big thrust
			game.sm.play(game.sm.sThrust)
			game.ship.isSliding = false
		}
		if rl.IsKeyDown('W') {
			game.ship.thrust(2.0)
		}
		if rl.IsKeyReleased('W') { // -----
			game.ship.thrust(0)
			game.ship.isSliding = true
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyDown('D') { // rotate right
			game.ship.rotate(.2)
		}
		if rl.IsKeyPressed(rl.KeyLeftControl) { // fire
			if game.missilesNo < maxMissiles {
				launchMissile(game)
			}

		}
		if rl.IsKeyDown(rl.KeyCapsLock) { // slow down rotation
			game.ship.m.rotSpeed *= 0.9
		}

		game.drawAndUpdate()
		dx, dy := rl.GetMouseDelta().X, rl.GetMouseDelta().X

		if !cursorEnabled && dx*dx+dy*dy > 16 {
			rl.EnableCursor()
			cursorEnabled = true
		}

	}
	game.finalize()
}
