package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ms rl.Music

func main() {
	// By default, Go programs run with GOMAXPROCS set to the number
	// of cores available; in prior releases it defaulted to 1.
	// Starting from Go 1.5, the default value is the number of cores.
	// You only need to explicitly set it if you are not okay with this
	// in newer Go versions.
	//runtime.GOMAXPROCS(8)

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Println("panic occurred:", err)
	// 	}
	// }()
	rl.SetTraceLog(rl.LogAll)
	// GetWindowHandle();
	game := newGame(1440, 720)

	for !rl.WindowShouldClose() {

		game.processKeys()
		game.processMouse()
		game.playMessages()
		game.drawAndUpdate()

	}

	rl.UnloadMusicStream(ms)
	game.finalize()
}

func (g *game) processKeys() {
	if rl.IsKeyPressed('Q') {
		if g.ship.cash > 16 {
			g.sm.play(sMissilesDlvrd)
			g.addParticle(newTextPart(g.ship.pos, g.ship.speed.MulA(0.5),
				"+20 missiles", 20, 1, 1, true, rl.Purple, rl.DarkPurple))
			g.ship.cash -= 16
			g.ship.missiles += 20
		}
	}
	if rl.IsKeyPressed('M') {
		if !g.sm.mute {
			g.sm.stopAll()
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
	if rl.IsKeyDown('S') { // hold thrust
		g.ship.thrust(0.5)
	}
	if rl.IsKeyReleased('S') { // end thrust
		g.ship.thrust(0)
		g.ship.isSliding = true
		g.sm.stop(sThrust)
	}
	if rl.IsKeyPressed(';') { //forceField

		g.ship.forceField = true
		g.sm.play(sForceField)
	}
	if rl.IsKeyDown(';') {
		if g.ship.energy > 0 {
			g.ship.energy -= 0.1
		}
	}
	if rl.IsKeyReleased(';') { // hold thrust
		g.ship.forceField = false
		g.sm.stop(sForceField)
	}
	if rl.IsKeyPressed('L') {
		g.weapon++
		g.weapon %= maxWeapons
		str := fmt.Sprintf(">%s<", weapons[g.weapon])
		g.addParticle(newTextPart(g.ship.pos, g.ship.speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))
	} // cycle weapon left
	if rl.IsKeyPressed('\'') {
		g.weapon++
		g.weapon %= maxWeapons
		str := fmt.Sprintf(">%s<", weapons[g.weapon])
		g.addParticle(newTextPart(g.ship.pos, g.ship.speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))

	} // cycle weapon right

	if rl.IsKeyPressed(rl.KeyF1) { /// debug
		debug = !debug
	}
	if rl.IsKeyPressed('R') { // reset shields
		g.sm.play(sOinx)
		g.ship.pos = V2{g.gW / 2, g.gH / 2}
		g.ship.speed = V2{0, 0}
		g.ship.shields = 100
		g.ship.energy = 1000
		g.ship.destroyed = false
	}
	if rl.IsKeyPressed('F') { // reset shields
		if g.ship.energy > 130 && g.ship.shields+13 < 100 {
			g.addParticle(newTextPart(g.ship.pos, g.ship.speed.MulA(0.5),
				"shields +13", 20, 1, 0.5, true, rl.Yellow, rl.Gold))
			g.sm.play(sChargeUp)
			g.ship.shields += 13
			g.ship.energy -= 130
		}
	}
	if rl.IsKeyPressed('P') { // pause
		g.paused = !g.paused
		g.sm.stopAll()
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

			if len(g.missiles) < maxMissiles {
				launchMissile(g)
				g.sm.play(sLaunch)
			}
		}
	}
	if rl.IsKeyDown(rl.KeyTab) { // slow down rotation
		g.ship.rotSpeed *= 0.9
	}

}
