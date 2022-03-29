package main

import (
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
	// GetWindowHandle();
	game := newGame(1440, 720)

	initMouse()

	for !rl.WindowShouldClose() {

		if !game.sm.isPlaying(sSpace) {
			game.sm.play(sSpace)

		}
		if !game.sm.isPlaying(sScore) {
			game.sm.play(sScore)
		}

		processKeys(game)
		processMouse()
		game.playMessages()
		game.drawAndUpdate()

	}
	game.finalize()
}

func processKeys(g *game) {
	if rl.IsKeyPressed('Q') {
		g.sm.playM(sMissilesDlvrd)
		if g.ship.cash > 16 {
			g.ship.cash -= 16
			g.ship.missiles += 10
		}
	}
	if rl.IsKeyPressed('M') {
		if !g.sm.mute {
			g.sm.stop(0)
		}
		g.sm.mute = !g.sm.mute

	}
	if rl.IsKeyDown('A') { // rotate left
		g.ship.rotate(-.2)
	}
	if rl.IsKeyPressed('S') { // small thrust
		g.sm.play(sThrust)
		g.ship.isSliding = false
	}
	if rl.IsKeyDown('S') {
		g.ship.thrust(0.5)
	}
	if rl.IsKeyReleased('S') { // -----
		g.ship.thrust(0)
		g.ship.isSliding = true
		g.sm.stop(sThrust)
	}
	if rl.IsKeyPressed('R') { // reset shields
		g.sm.play(sOinx)
		g.ship.shields = 100
		g.ship.destroyed = false
	}
	if rl.IsKeyPressed('W') { // big thrust
		g.sm.play(sThrust)
		g.ship.isSliding = false
	}
	if rl.IsKeyDown('W') {
		g.ship.thrust(1.0)
	}
	if rl.IsKeyReleased('W') { // -----
		g.ship.thrust(0)
		g.ship.isSliding = true
		g.sm.stop(sThrust)
	}
	if rl.IsKeyDown('D') { // rotate right
		g.ship.rotate(.2)
	}
	if rl.IsKeyPressed(rl.KeyLeftControl) { // fire
		if g.ship.missiles > 0 {
			g.ship.missiles--

			if g.missilesNo < maxMissiles {
				launchMissile(g)
				g.sm.playM(sLaunch)
			}
		}

	}
	if rl.IsKeyDown(rl.KeyTab) { // slow down rotation
		g.ship.m.rotSpeed *= 0.9
	}
}

var cursorEnabled bool

func initMouse() {
	rl.DisableCursor()
	cursorEnabled = false
}
func processMouse() {
	dx, dy := rl.GetMouseDelta().X, rl.GetMouseDelta().X

	if !cursorEnabled && dx*dx+dy*dy > 16 {
		rl.EnableCursor()
		cursorEnabled = true
	}

}
