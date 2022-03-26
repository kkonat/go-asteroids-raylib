package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// By default, Go programs run with GOMAXPROCS set to the number
	// of cores available; in prior releases it defaulted to 1.
	// Starting from Go 1.5, the default value is the number of cores.
	// You only need to explicitly set it if you are not okay with this
	// in newer Go versions.
	//runtime.GOMAXPROCS(8)

	rl.SetTraceLog(rl.LogAll)

	rand.Seed(time.Now().UnixNano())
	_initNoise()
	game := newGame(1440, 720)
	rl.DisableCursor()
	cursorEnabled := false

	for !rl.WindowShouldClose() {

		if !game.sm.isPlaying(0) {
			game.sm.play(sSpace)
			//			fmt.Println("start backg sound")
		}

		if rl.IsKeyPressed('Q') {
			game.sm.playM(sOinx)
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
			game.sm.play(sThrust)
			game.ship.isSliding = false
		}
		if rl.IsKeyDown('S') {
			game.ship.thrust(1.0)
		}
		if rl.IsKeyReleased('S') { // -----
			game.ship.thrust(0)
			game.ship.isSliding = true
			game.sm.stop(sThrust)
		}
		if rl.IsKeyPressed('W') { // big thrust
			game.sm.play(sThrust)
			game.ship.isSliding = false
		}
		if rl.IsKeyDown('W') {
			game.ship.thrust(2.0)
		}
		if rl.IsKeyReleased('W') { // -----
			game.ship.thrust(0)
			game.ship.isSliding = true
			game.sm.stop(sThrust)
		}
		if rl.IsKeyDown('D') { // rotate right
			game.ship.rotate(.2)
		}
		if rl.IsKeyPressed(rl.KeyLeftControl) { // fire
			if game.missilesNo < maxMissiles {
				launchMissile(game)
				game.sm.playM(sLaunch)
			}

		}
		if rl.IsKeyDown(rl.KeyCapsLock) { // slow down rotation
			game.ship.m.rotSpeed *= 0.9
		}

		game.drawAndUpdate()
		//printMemoryUsage()
		dx, dy := rl.GetMouseDelta().X, rl.GetMouseDelta().X

		if !cursorEnabled && dx*dx+dy*dy > 16 {
			rl.EnableCursor()
			cursorEnabled = true
		}

	}
	game.finalize()
}
