package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	caption      = "test bum bum game"
	rSpeedMax    = 1
	initialRocks = 12
	maxRocks     = 100
	maxMissiles  = 50
)

type game struct {
	sm         *soundManager
	sf         *starfield
	ship       *ship
	rocks      [maxRocks]*Rock
	rocksNo    int
	missiles   [maxMissiles]*missile
	missilesNo int
	sW, sH     int32
	gW, gH     float64
}

var tnow, tprev int64

func newGame(w, h int32) *game {

	const safeCircle = 200

	g := new(game)
	g.sW, g.sH = w, h
	g.gW, g.gH = float64(w), float64(h)

	cX, cY := float64(w/2), float64(h/2)

	g.sf = newStarfield(w, h)
	g.sm = newSoundManager(true)
	g.ship = newShip(cX, cY, 1000, 1000)

	i := 0
	for i < initialRocks { // ( cx +r )  ( nr.x +nr.r)
		nr := newRockRandom(g)
		if cX+safeCircle < nr.m.pos.x+nr.radius || cX-safeCircle > nr.m.pos.x-nr.radius ||
			cY+safeCircle < nr.m.pos.y+nr.radius || cY-safeCircle > nr.m.pos.y-nr.radius {
			g.rocks[i] = nr
			i++
		}
	}
	g.rocksNo = i

	g.prepareDisplay()
	tnow = time.Now().Local().UnixMicro()
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
func (gme *game) moveRocks(dt float64) {

	for i := 0; i < gme.rocksNo; i++ { // move rocks
		gme.rocks[i].m.Move(dt)
	}
}
func (gme *game) moveMissiles(dt float64) {

	for i := 0; i < gme.rocksNo; i++ { // move rocks
		gme.rocks[i].m.Move(dt)
	}

	for i := 0; i < gme.missilesNo; i++ { // move missiles
		gme.missiles[i].m.Move(dt)

	}
}
func (gme *game) drawAndUpdate() {

	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	gme.sf.draw() // draw starfield

	for i := 0; i < gme.rocksNo; i++ { // draw rocks
		gme.rocks[i].Draw()
	}

	for i := 0; i < gme.missilesNo; i++ { // draw missiles
		gme.missiles[i].Draw()

	}

	gme.ship.Draw() // draw ship

	gme.drawStatusBar() // draw status bar on top of everything
	gme.sm.doFade()     // fade out sounds if needed

	rl.EndDrawing()

	tnow = time.Now().UnixMicro()
	elapsed := tnow - tprev
	tprev = tnow
	dt := float64(elapsed) / 16666.0
	fmt.Println(dt)
	gme.ship.m.Move(dt)
	gme.moveRocks(dt)
	gme.moveMissiles(dt)
	gme.process_missile_hits()
	gme.constrainShip()
	gme.constrainRocks()
	gme.constrainMissiles()
	//	gme.animatestuff()

}

func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.unloadAll()
}
