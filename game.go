package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type game struct {
	sm     *soundManager
	sf     *starfield
	ship   *ship
	sW, sH int32
}

const (
	caption = "test bum bum game"
)

func newGame(w, h int32) *game {
	rand.Seed(time.Now().UnixNano())

	g := new(game)
	g.sW, g.sH = w, h

	g.sf = newStarfield(w, h)

	g.sm = newSoundManager()

	g.ship = newShip(1000, 1000)
	g.ship.pos = V2{720, 360}

	g.prepareDisplay()
	return g
}

func (g *game) drawStatusBar() {

	rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
	rl.DrawText(caption, 20, g.sH-20, 20, rl.Magenta)
	//rl.DrawLine(18, 42, g.sW-18, 42, rl.Black)
	rl.DrawFPS(g.sW-80, g.sH-20)
}
func (g *game) drawGrid() {
	const step = 40

	for i := int32(0); i < g.sH; i += step {
		rl.DrawLine(0, i, g.sW, i, rl.DarkGray)
	}
	for i := int32(0); i < g.sW; i += step {
		rl.DrawLine(i, 0, i, g.sH, rl.DarkGray)
	}

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
	//gme.drawGrid()
	gme.sf.draw()
	gme.ship.Draw()
	gme.drawStatusBar() // draw on top of everything

	rl.EndDrawing()

}
func (g *game) resizeDisplay() {
	g.sW = int32(rl.GetScreenWidth())
	g.sH = int32(rl.GetScreenHeight())
}
func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.unloadAll()
}
