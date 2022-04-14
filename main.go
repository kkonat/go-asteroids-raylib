package main

import (
	"fmt"
	"math/rand"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ms rl.Music
var Game *game

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
	rl.SetTraceLog(rl.LogNone)

	rand.Seed(time.Now().UnixNano())
	_initNoise()

	Game = newGame(1440, 720)
	defer func() {
		rl.UnloadMusicStream(ms)
		Game.finalize()
	}()

	gui.LoadGuiStyle("res/dark.style")

	for !rl.WindowShouldClose() {
		Game.processKeys()
		Game.processMouse()
		Game.playMessages()
		Game.GameDraw()
		Game.GameUpdate()
	}
}

var fflight *OmniLight

func (g *game) processKeys() {
	if rl.IsKeyPressed('Q') {
		wpn := Game.weapons[Game.curWeapon]
		cost := int(20 * wpn.cost)
		if Game.ship.cash > cost {
			Game.sm.Play(sMissilesDlvrd)
			Game.addParticle(newTextPart(Game.ship.pos, Game.ship.speed.MulA(0.5),
				"+ 20 x "+wpn.name, 20, 1, 1, true, rl.Purple, rl.DarkPurple))
			Game.ship.cash -= cost
			wpn.curCap += 20
			Game.weapons[Game.curWeapon] = wpn
		}
	}
	if rl.IsKeyPressed('M') {
		if !Game.sm.Mute {
			Game.sm.StopAll()
		}
		Game.sm.Mute = !Game.sm.Mute

	}
	if rl.IsKeyDown('A') { // rotate left
		Game.ship.rotate(-.2)
	}
	if rl.IsKeyPressed('S') { // small thrust
		Game.sm.Play(sThrust)
		Game.ship.isSliding = false
	}
	if rl.IsKeyDown('S') { // hold thrust
		Game.ship.thrust(0.5)
	}
	if rl.IsKeyReleased('S') { // end thrust
		Game.ship.thrust(0)
		Game.ship.isSliding = true
		Game.sm.Stop(sThrust)
	}
	if rl.IsKeyPressed(';') { //forceField
		fflight = &OmniLight{Game.ship.pos, Color{0, 0.78, 0.78, 1.0}, 100}
		Game.VisibleLights.AddLight(fflight)
		Game.ship.forceField = true
		Game.sm.Play(sForceField)
	}
	if rl.IsKeyDown(';') {
		fflight.Pos = Game.ship.pos
		if Game.ship.energy > 0 {
			Game.ship.energy -= 0.1
		}
	}
	if rl.IsKeyReleased(';') { // hold thrust
		Game.VisibleLights.DeleteLight(fflight)
		Game.ship.forceField = false
		Game.sm.Stop(sForceField)
	}
	if rl.IsKeyPressed('L') {
		if Game.curWeapon > 0 {
			Game.curWeapon--
		} else {
			Game.curWeapon = len(Game.weapons) - 1
		}

		str := fmt.Sprintf(">%s<", (Game.weapons)[Game.curWeapon].name)
		Game.addParticle(newTextPart(Game.ship.pos, Game.ship.speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))
	} // cycle weapon left
	if rl.IsKeyPressed('\'') {
		Game.curWeapon++
		Game.curWeapon %= len(Game.weapons)
		str := fmt.Sprintf(">%s<", Game.weapons[Game.curWeapon].name)
		Game.addParticle(newTextPart(Game.ship.pos, Game.ship.speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))

	} // cycle weapon right

	if rl.IsKeyPressed(rl.KeyF1) { /// debug
		debug = !debug
	}
	if rl.IsKeyPressed(rl.KeyF2) { /// debug
		showgui = !showgui
	}
	if rl.IsKeyPressed('R') { // reset shields
		Game.sm.Play(sOinx)
		Game.ship.pos = V2{X: Game.gW / 2, Y: Game.gH / 2}
		Game.ship.speed = V2{}
		Game.ship.shields = 100
		Game.ship.energy = 1000
		Game.ship.destroyed = false
	}
	if rl.IsKeyPressed('F') { // reset shields
		if Game.ship.energy > 130 && Game.ship.shields+13 < 100 {
			Game.addParticle(newTextPart(Game.ship.pos, Game.ship.speed.MulA(0.5),
				"shields +13", 20, 1, 0.5, true, rl.Yellow, rl.Gold))
			Game.sm.Play(sChargeUp)
			Game.ship.shields += 13
			Game.ship.energy -= 130
		}
	}
	if rl.IsKeyPressed('P') { // pause
		Game.paused = !Game.paused
		Game.sm.StopAll()
	}
	if rl.IsKeyPressed('W') { // big thrust
		Game.sm.Play(sThrust)
		Game.ship.isSliding = false
	}
	if rl.IsKeyDown('W') {
		Game.ship.thrust(1.0)
	}
	if rl.IsKeyReleased('W') { // -----
		Game.ship.thrust(0)
		Game.ship.isSliding = true
		Game.sm.Stop(sThrust)
	}
	if rl.IsKeyDown('D') { // rotate right
		Game.ship.rotate(.2)
	}
	if rl.IsKeyPressed(rl.KeyLeftControl) { // fire
		wpn := Game.weapons[Game.curWeapon]
		if wpn.curCap > 0 {
			wpn.curCap -= 1
			Game.weapons[Game.curWeapon] = wpn
			if len(Game.missiles) < maxMissiles {
				launchMissile(g, Game.curWeapon)
				Game.sm.Play(sLaunch)
			}
		}
	}
	if rl.IsKeyDown(rl.KeyTab) { // slow down rotation
		Game.ship.rotSpeed *= 0.9
	}

}
