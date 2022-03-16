package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type game struct {
	sm       *soundManager
	sf       *starfield
	ship     *Ship
	rocks    []Rock
	missiles []missile
	sW, sH   int32
}

const (
	caption        = "test bum bum game"
	noRocks        = 12
	RockSpeedrange = 1
)

func newGame(w, h int32) *game {
	rand.Seed(time.Now().UnixNano())

	g := new(game)
	g.sW, g.sH = w, h

	g.sf = newStarfield(w, h)

	g.sm = newSoundManager()

	g.ship = newShip(720, 360, 1000, 1000)
	for i := 0; i < noRocks; i++ {
		g.rocks = append(g.rocks,
			*newRock(V2{rand.Float32() * float32(g.sW), rand.Float32() * float32(g.sH)},
				V2{rand.Float32()*RockSpeedrange*2.0 - RockSpeedrange, rand.Float32()*RockSpeedrange*2.0 - RockSpeedrange},
				rand.Float32()*0.2-0.1))
	}

	g.prepareDisplay()
	return g
}
func (g *game) constrainShip() {
	const limit = 40
	const getback = 0.5
	if g.ship.shape.pos.x < limit || g.ship.shape.pos.x > float32(g.sW-limit) || g.ship.shape.pos.y < limit || g.ship.shape.pos.y > float32(g.sH-limit) {
		if g.ship.Slides {
			g.ship.shape.speed = V2V(g.ship.shape.speed, 0.9)
			if g.ship.shape.speed.x*g.ship.shape.speed.x+g.ship.shape.speed.x*g.ship.shape.speed.x < 0.01 {
				g.ship.Slides = false
			}
		} else {
			if g.ship.shape.pos.x < float32(limit) {
				g.ship.shape.speed.x = getback
			}
			if g.ship.shape.pos.x > float32(g.sW-limit) {
				g.ship.shape.speed.x = -getback
			}
			if g.ship.shape.pos.y < float32(limit) {
				g.ship.shape.speed.y = getback
			}
			if g.ship.shape.pos.y > float32(g.sH-limit) {
				g.ship.shape.speed.y = -getback
			}
		}

	}
}
func (g *game) constrainRocks() {
	const limit = -50
	for i := range g.rocks {
		p := &g.rocks[i].shape.pos
		s := &g.rocks[i].shape.speed
		r := g.rocks[i].radius
		if (p.x+r < limit && s.x < 0) || (p.x-r > float32(g.sW)-limit && s.x > 0) ||
			(p.y+r < limit && s.y < 0) || (p.y-r > float32(g.sH)-limit && s.y > 0) {
			g.rocks[i].Randomize()
			g.rocks[i].shape.speed = V2{rand.Float32()*RockSpeedrange - RockSpeedrange/2, rand.Float32()*RockSpeedrange - RockSpeedrange/2}
			g.rocks[i].shape.pos.x = float32(rand.Int31n(g.sW))
			g.rocks[i].shape.pos.y = float32(rand.Int31n(g.sH))

			sect := rand.Intn(4)

			switch sect {
			case 0:
				{
					p.x = -r + limit
					if s.x < 0 {
						s.x = -s.x
					}
				}
			case 1:
				{
					p.y = -r + limit
					if s.y < 0 {
						s.x = -s.x
					}
				}
			case 2:
				{
					p.x = float32(g.sW) + r - limit
					if s.x > 0 {
						s.x = -s.x
					}
				}
			case 3:
				{
					p.y = float32(g.sH) + r - limit
					if s.y > 0 {
						s.y = -s.y
					}
				}
			} // respawn rock in a new sector

		}
	}
}
func (g *game) deleteMissile(i int) {
	g.missiles[i] = g.missiles[len(g.missiles)-1]
	g.missiles = g.missiles[:len(g.missiles)-1]
}

func (g *game) constrainMissiles() {

	const limit = -10
	for i := range g.missiles {
		p := g.missiles[i].shape.pos
		if p.x <= limit || p.x > float32(g.sW-limit) ||
			p.y <= limit || p.y > float32(g.sH-limit) {
			g.deleteMissile(i)

			break
		}
	}
}
func (g *game) drawStatusBar() {

	rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
	rl.DrawText(caption, 20, g.sH-20, 20, rl.Magenta)
	//rl.DrawLine(18, 42, g.sW-18, 42, rl.Black)
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
	// draw starfield
	gme.sf.draw()

	// draw ship
	gme.ship.Draw()

	for i := range gme.rocks {
		gme.rocks[i].Draw()
	}

	for i := range gme.missiles {
		gme.missiles[i].Draw()
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
