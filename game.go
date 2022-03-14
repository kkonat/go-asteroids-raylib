package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type game struct {
	sW, sH int32
}

const (
	caption = "test bum bum game"
)

func newGame(w, h int32) *game {
	g := new(game)
	g.sW, g.sH = w, h
	g.prepareDisplay()
	return g
}

func (g *game) displayStatus() {

	rl.DrawRectangleV(rl.NewVector2(0, 20), rl.NewVector2(float32(g.sW), 26), rl.White)
	rl.DrawText(caption, 20, 20, 20, rl.DarkGray)
	rl.DrawLine(18, 42, g.sW-18, 42, rl.Black)
	rl.DrawFPS(g.sW-80, 20)
}
func (g *game) prepareDisplay() {

	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)

	rl.InitWindow(g.sW, g.sH, caption)
	rl.MaximizeWindow()

	//rl.SetTargetFPS(60)
}

func (g *game) drawGame(w *world) {

	rl.UpdateCamera(&w.camera)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.BeginMode3D(w.camera)

	w.drawBackground()
	w.drawObjects()

	rl.EndMode3D()

	g.displayStatus()
	//rl.DrawTexture(tex, 0, 0, rl.White)
	rl.EndDrawing()

}
func (g *game) resizeDisplay() {
	g.sW = int32(rl.GetScreenWidth())
	g.sH = int32(rl.GetScreenHeight())
}
func (g *game) finalize() {
	rl.CloseWindow()
}
