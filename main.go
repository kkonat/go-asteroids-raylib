package main

import (
	"fmt"
	"math/rand"
	"time"

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
	//rl.SetTraceLog(rl.LogAll)

	rand.Seed(time.Now().UnixNano())
	_initNoise()

	Game = newGame(1440, 720)
	defer func() {
		rl.UnloadMusicStream(ms)
		Game.finalize()
	}()

	// gui.LoadGuiStyle("res/dark.style")

	var tnow, tprev int64
	for !rl.WindowShouldClose() {
		processKeys()
		processMouse()

		Game.playCyclicMessages()
		Game.GameDraw()

		tnow = time.Now().UnixMicro()
		elapsed := tnow - tprev
		tprev = tnow
		dt := float64(elapsed) / 16666.0

		Game.GameUpdate(dt)
	}
}

var fflight *OmniLight

func shipNavKeys() {
	// rotate ship left
	if rl.IsKeyDown('A') {
		Game.ship.rotate(-.2)
	}
	// rotate right
	if rl.IsKeyDown('D') {
		Game.ship.rotate(.2)
	}
	// start applying  thrust
	if rl.IsKeyPressed('W') || rl.IsKeyPressed('S') {
		Game.sm.Play(sThrust)
		Game.ship.isSliding = false
	}
	// continue applying more thrust
	if rl.IsKeyDown('W') {
		Game.ship.Thrust(1.0)
	}
	// continue applying small thrust
	if rl.IsKeyDown('S') {
		Game.ship.Thrust(0.5)
	}
	// stop applying thrust
	if rl.IsKeyReleased('W') || rl.IsKeyReleased('S') {
		Game.ship.Thrust(0)
		Game.ship.isSliding = true
		Game.sm.Stop(sThrust)
	}
	// slow down rotation
	if rl.IsKeyDown(rl.KeyTab) {
		Game.ship.RotSpeed *= 0.9
	}

}
func shipSystemsKeys() {
	if rl.IsKeyPressed('Z') {
		Game.ship.SpotlightMode()
	}

	// start emit forceField
	if rl.IsKeyPressed(';') {
		fflight = &OmniLight{Game.ship.Pos, Color{0, 0.78, 0.78, 1.0}, 100}
		Game.VisibleLights.AddLight(fflight)
		Game.ship.forceField = true
		Game.sm.Play(sForceField)
	}
	// continue emit forceField
	if rl.IsKeyDown(';') {
		fflight.Pos = Game.ship.Pos
		if Game.ship.energy > 0 {
			Game.ship.energy -= 0.1
		}
	}
	// stop emit forceField
	if rl.IsKeyReleased(';') {
		Game.VisibleLights.DeleteLight(fflight)
		Game.ship.forceField = false
		Game.sm.Stop(sForceField)
	}
	// "shields +13"
	if rl.IsKeyPressed('F') {
		if Game.ship.energy > 130 && Game.ship.shields+13 < 100 {
			Game.addParticle(newTextPart(Game.ship.Pos, Game.ship.Speed.MulA(0.5),
				"shields +13", 20, 1, 0.5, true, rl.Yellow, rl.Gold))
			Game.sm.Play(sChargeUp)
			Game.ship.shields += 13
			Game.ship.energy -= 130
		}
	}
}
func weaponsKeys() {
	// Buy missiles
	if rl.IsKeyPressed('Q') {
		wpn := Game.weapons[Game.curWeapon]
		cost := int(20 * wpn.cost)
		if Game.ship.cash >= cost {
			Game.sm.Play(sMissilesDlvrd)
			Game.addParticle(newTextPart(Game.ship.Pos, Game.ship.Speed.MulA(0.5),
				"+ 20 x "+wpn.name, 20, 1, 1, true, rl.Purple, rl.DarkPurple))
			Game.ship.cash -= cost
			wpn.curCap += 20
			Game.weapons[Game.curWeapon] = wpn
		}
	}
	// cycle weapon -
	if rl.IsKeyPressed('L') {
		if Game.curWeapon > 0 {
			Game.curWeapon--
		} else {
			Game.curWeapon = len(Game.weapons) - 1
		}

		str := fmt.Sprintf(">%s<", (Game.weapons)[Game.curWeapon].name)
		Game.addParticle(newTextPart(Game.ship.Pos, Game.ship.Speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))
	}
	// cycle weapon +
	if rl.IsKeyPressed('\'') {
		Game.curWeapon++
		Game.curWeapon %= len(Game.weapons)
		str := fmt.Sprintf(">%s<", Game.weapons[Game.curWeapon].name)
		Game.addParticle(newTextPart(Game.ship.Pos, Game.ship.Speed.MulA(0.5),
			str, 20, 1, 1, true, rl.Purple, rl.Red))
	}
	// fire
	if rl.IsKeyPressed(rl.KeyLeftControl) {
		wpn := Game.weapons[Game.curWeapon]
		if wpn.curCap > 0 {
			wpn.curCap -= 1
			Game.weapons[Game.curWeapon] = wpn
			if len(Game.missiles) < maxMissiles {
				Game.launchMissile()
				Game.sm.Play(sLaunch)
			}
		}
	}
}

func gameKeys() {
	// mute / unmute
	if rl.IsKeyPressed('M') {
		if !Game.sm.Mute {
			Game.sm.StopAll()
		}
		Game.sm.Mute = !Game.sm.Mute
	}
	// debug
	if rl.IsKeyPressed(rl.KeyF3) {
		debug = !debug
	}
	// gui
	if rl.IsKeyPressed(rl.KeyF2) {
		showgui = !showgui
	}
	if rl.IsKeyPressed(rl.KeyF1) {
		showKeys = !showKeys
	}
	// reset shields
	if rl.IsKeyPressed('R') {
		Game.sm.Play(sOinx)
		Game.ship.Respawn()
	}
	// pause
	if rl.IsKeyPressed('P') {
		Game.paused = !Game.paused
		Game.sm.StopAll()
	}
}
func processKeys() {
	shipNavKeys()
	weaponsKeys()
	shipSystemsKeys()
	gameKeys()
}
