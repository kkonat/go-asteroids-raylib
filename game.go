package main

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	// "github.com/golang-ui/nuklear/nk"	// interface, maybe
	"golang.org/x/exp/constraints"
)

const (
	caption          = "test boom boom game"
	rSpeedMax        = 1
	noPreferredRocks = 30
	maxRocks         = 100
	maxMissiles      = 50
	maxParticles     = 50
	FPS              = 60
	shieldsLowLimit  = 25
	ammoLowLimit     = 20
)

type game struct {
	sm   *soundManager
	sprm *spriteManager
	sf   *starfield
	time []float32 // gsls uniform [1]float32
	ship *ship

	rocks     []*Rock
	missiles  []*missile
	particles []particle

	qt *QuadTree

	//objects     []*object

	sW, sH        int32
	gW, gH        float64
	ufo           rl.Texture2D
	paused        bool
	cursorEnabled bool
	debug         bool
}

// global time counters
var tnow, tprev, mAmmoLastPlayed, mShieldsLastPlayed int64
var tickTock uint8

func newGame(w, h int32) *game {

	rand.Seed(time.Now().UnixNano())
	_initNoise()

	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)

	rl.InitWindow(w, h, caption)
	// handle := rl.GetWindowHandle()
	// println(handle)
	// ctx := nk.NkPlatformInit(handle, nk.PlatformInstallCallbacks)
	// println(ctx)
	// rl.MaximizeWindow()

	rl.SetTargetFPS(FPS)

	g := new(game)

	g.qt = NewQuadTree(0, Rect{0, 0, w, h})
	g.initMouse()
	g.paused = false
	g.sW, g.sH = w, h
	g.gW, g.gH = float64(w), float64(h)

	g.time = make([]float32, 1)
	g.sf = newStarfield(w, h, g.time)

	g.sm = newSoundManager(false)
	g.sprm = newSpriteManager()

	//g.objects = make([]*object, 0, 20)

	g.ship = newShip(float64(w/2), float64(h/2), 1000, 1000)

	generateRocks(g, noPreferredRocks)

	tprev = time.Now().Local().UnixMicro()
	return g
}

func (g game) playMessages() {
	if g.ship.missiles < ammoLowLimit {
		t := time.Now().Local().Unix()
		if mAmmoLastPlayed == 0 || t-mAmmoLastPlayed > 15 {
			{
				g.sm.playM(sAmmoLow)
				mAmmoLastPlayed = time.Now().Local().Unix()
			}
		}
	} else {
		mAmmoLastPlayed = 0
	}
	if g.ship.shields < shieldsLowLimit {
		t := time.Now().Local().Unix()
		if mShieldsLastPlayed == 0 || t-mShieldsLastPlayed > 17 {
			{
				g.sm.playM(sShieldsLow)
				mShieldsLastPlayed = time.Now().Local().Unix()
			}
		}
	} else {
		mShieldsLastPlayed = 0
	}
}

func flashColor[T constraints.Ordered](col rl.Color, warn, low, val T) rl.Color {
	if val < low {
		if tickTock%20 > 10 {
			return rl.Red
		} else {
			return color.RGBA{127, 0, 0, 255}
		}
	} else if val < warn {
		return rl.Beige
	} else {
		return col
	}
}

func (g *game) drawStatusBar() {
	if !g.paused {
		rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
		_multicolorText(20, g.sH-20, 20,
			"Cash:", rl.Purple, g.ship.cash, rl.Purple, 30, 10,
			"Shields:", rl.Purple, int(g.ship.shields), flashColor(rl.Purple, 50, shieldsLowLimit, int(g.ship.shields)),
			"Energy:", rl.Purple, int(g.ship.energy), flashColor(rl.Purple, 500, 100, int(g.ship.energy)),
			"Missiles:", rl.Purple, g.ship.missiles, flashColor(rl.Purple, 30, ammoLowLimit, g.ship.missiles))

		rl.DrawFPS(g.sW-80, g.sH-20)
	} else {
		var col rl.Color
		if tickTock%20 > 10 {
			col = rl.DarkPurple
		} else {
			col = rl.Purple
		}
		rl.DrawText("**** GAME PAUSED ***", 20, g.sH-20, 20, col)
	}
	tickTock++
}

func (gme *game) addParticle(p particle) {
	if len(gme.particles) < maxParticles {
		gme.particles = append(gme.particles, p)
	}
}
func (gme *game) animateParticles() {
	var i int
	for i < len(gme.particles) {
		gme.particles[i].Animate()
		if gme.particles[i].canDelete() {
			gme.particles = append(gme.particles[:i], gme.particles[i+1:]...)
		} else {
			i++
		}
	}
}
func (gme *game) drawRocks() {
	for _, r := range gme.rocks { // draw rocks
		r.Draw()
	}
}
func (gme *game) drawMissiles() {
	for _, m := range gme.missiles {
		m.Draw()
	}
}
func (gme *game) drawParticles() {
	for i := 0; i < len(gme.particles); i++ {
		gme.particles[i].Draw()
	}
}

func (gme *game) moveRocks(dt float64) {
	for i := 0; i < len(gme.rocks); i++ { // move rocks
		go gme.rocks[i].m.Move(dt)
	}
	wg.Done()
}
func (gme *game) moveMissiles(dt float64) {
	for i := range gme.missiles { // move missiles
		go gme.missiles[i].m.Move(dt)
	}
	wg.Done()
}

var wg sync.WaitGroup

func (gme *game) drawAndUpdate() {

	if !gme.paused {
		if !gme.sm.isPlaying(sSpace) {
			gme.sm.play(sSpace)

		}
		if !gme.sm.isPlaying(sScore) {
			gme.sm.play(sScore)
		}
	}
	// DRAWING ---------------------------------------------------------------------
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)
	gme.sf.draw() // draw starfield

	gme.drawRocks()

	gme.debugQt()

	gme.drawMissiles()
	gme.drawParticles()
	gme.ship.Draw()
	gme.drawStatusBar()

	rl.EndDrawing()

	gme.sm.doFade() // fade out sounds if needed

	// UPDATING ---------------------------------------------------------------------
	tnow = time.Now().UnixMicro()
	elapsed := tnow - tprev
	tprev = tnow
	dt := float64(elapsed) / 16666.0

	if !gme.paused {

		gme.time[0] += 0.01 // glsl uniform for starfield shader

		gme.ship.m.Move(dt)

		gme.ship.chargeUp() // chargeup ship

		wg.Add(1) // Waitgroup
		gme.moveRocks(dt)
		wg.Add(1)
		gme.moveMissiles(dt)

		wg.Wait()

		// t0 := time.Now().UnixNano()
		gme.process_missile_hits()
		gme.process_ship_hits()
		// t0 = time.Now().UnixNano() - t0
		// var tmax int64
		// comps := gme.missilesNo * gme.rocksNo
		// if tmax < t0 {
		// 	tmax = t0
		// 	fmt.Printf("[%d comps took %.3f us]\n", comps, float64(tmax)/1000)
		// }
		gme.constrainShip()

		go gme.constrainRocks()
		go gme.constrainMissiles()
		go gme.animateParticles()
	}
}
func (g *game) processKeys() {
	if rl.IsKeyPressed('Q') {
		if g.ship.cash > 16 {
			g.sm.playM(sMissilesDlvrd)
			g.addParticle(newTextPart(g.ship.m.pos, g.ship.m.speed.MulA(0.5), "+20 missiles", 20, 3, rl.Purple, rl.DarkPurple))
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
	if rl.IsKeyPressed(rl.KeyF1) { /// debug
		g.debug = !g.debug
	}
	if rl.IsKeyPressed('R') { // reset shields
		g.sm.play(sOinx)
		g.ship.m.pos = V2{g.gW / 2, g.gH / 2}
		g.ship.m.speed = V2{0, 0}
		g.ship.shields = 100
		g.ship.energy = 1000
		g.ship.destroyed = false
	}
	if rl.IsKeyPressed('F') { // reset shields
		if g.ship.energy > 130 && g.ship.shields+13 < 100 {
			g.addParticle(newTextPart(g.ship.m.pos, g.ship.m.speed.MulA(0.5), "shields +13", 20, 3, rl.Yellow, rl.Gold))
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
				g.sm.playM(sLaunch)
			}
		}
	}
	if rl.IsKeyDown(rl.KeyTab) { // slow down rotation
		g.ship.m.rotSpeed *= 0.9
	}

}

func (g *game) initMouse() {
	rl.DisableCursor()
	g.cursorEnabled = false
}
func (g *game) processMouse() {
	dx, dy := rl.GetMouseDelta().X, rl.GetMouseDelta().X

	if !g.cursorEnabled && dx*dx+dy*dy > 16 {
		rl.EnableCursor()
		g.cursorEnabled = true
	}
}
func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.unloadAll()
	g.sprm.unloadAll()
}

// -- debug
func drawQt(qt *QuadTree) {
	rl.DrawRectangleLines(qt.Bounds.x, qt.Bounds.y, qt.Bounds.w, qt.Bounds.h, rl.DarkGray)
	for i := 0; i < 4; i++ {
		if qt.Nodes[i] != nil {
			drawQt(qt.Nodes[i])
		}
	}
}
func (gme *game) debugQt() {
	if gme.debug {
		// str := fmt.Sprintf("[%d,%d]",int32(gme.ship.m.pos.x),int32(gme.ship.m.pos.y))
		// rl.DrawText(str, int32(gme.ship.m.pos.x),int32(gme.ship.m.pos.y), 20, rl.Gray)
		gme.qt.Clear()
		for _, r := range gme.rocks {
			gme.qt.Insert(newCircleV2(r.m.pos, r.radius))
		}

		largerCircle := newCircleV2(gme.ship.m.pos, 20)

		potCols := gme.qt.MayCollide(largerCircle)
		for _, c := range potCols {
			rl.DrawRectangleLines(c.rect.x, c.rect.y, c.rect.w, c.rect.h, rl.Beige)
		}
		drawQt(gme.qt)
		// for i, missile := range gme.missiles {
		// 	largerCircle = newCircleV2(missile.m.pos, 10)
		// 	potCols = gme.qt.MayCollide(largerCircle)
		// 	for _, c := range potCols {
		// 		rl.DrawRectangleLines(c.rect.x, c.rect.y, c.rect.w, c.rect.h, rl.DarkGreen)
		// 		str := fmt.Sprintf("[%d]", i)
		// 		rl.DrawText(str, c.rect.x+int32(i*16), c.rect.y, 16, rl.Lime)
		// 	}
		// }
	}
}
