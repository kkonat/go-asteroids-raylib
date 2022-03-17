package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type game struct {
	sm         *soundManager
	sf         *starfield
	ship       *Ship
	rocks      []Rock
	missiles   [50]*missile
	missilesNo int
	sW, sH     int32
}

const (
	caption   = "test bum bum game"
	rSpeedMax = 1
)

func newGame(w, h int32) *game {
	rand.Seed(time.Now().UnixNano())

	g := new(game)
	g.sW, g.sH = w, h

	g.sf = newStarfield(w, h)

	g.sm = newSoundManager(true)

	g.ship = newShip(720, 360, 1000, 1000)

	for i := 0; i < noRocks; i++ {
		nr := *newRock(g)
		g.rocks = append(g.rocks, nr)
	}

	g.prepareDisplay()
	return g
}

func (g *game) drawStatusBar() {

	rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
	rl.DrawText(caption, 20, g.sH-20, 20, rl.Magenta)
	rl.DrawFPS(g.sW-80, g.sH-20)
}

func (g *game) prepareDisplay() {

	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)

	rl.InitWindow(g.sW, g.sH, caption)

	rl.MaximizeWindow()

	rl.SetTargetFPS(60)
}

func (gme *game) drawGame() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	gme.sf.draw() // draw starfield

	gme.ship.Draw() // draw ship

	for i := range gme.rocks { // draw rocks
		gme.rocks[i].Draw()
	}

	i := 0
	for i < gme.missilesNo { // draw missiles
		gme.missiles[i].Draw()

		i++
	}

	gme.drawStatusBar() // draw status bar on top of everything
	gme.sm.doFade()     // fade out sounds if needed

	rl.EndDrawing()

	gme.constrainShip()
	gme.constrainRocks()
	gme.constrainMissiles()
}

func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.unloadAll()
}
