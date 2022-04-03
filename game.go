package main

import (
	"image/color"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	// "github.com/golang-ui/nuklear/nk"	// interface, maybe
	"golang.org/x/exp/constraints"
)

const (
	caption            = "test boom boom game"
	rSpeedMax          = 1
	noPreferredRocks   = 50
	PrefferredRockSize = 80
	maxRocks           = 100
	maxMissiles        = 50
	maxParticles       = 50
	FPS                = 60
	shieldsLowLimit    = 25
	ammoLowLimit       = 20

	startWithDebugOn = true
	startMuted       = true
)

var vectorFont rl.Font
var debug bool

// type rqtElem dlElem[*Rock]

// func (e rqtElem) bRect() Rect {
// 	return e.data.bRect()
// }

type RockListEl struct {
	ListEl[*Rock]
}

func (r RockListEl) bRect() Rect {
	x := int32(r.ListEl.Value.pos.x)
	y := int32(r.ListEl.Value.pos.y)
	rad := int32(r.ListEl.Value.radius)
	return Rect{x - rad, y - rad, rad * 2, rad * 2}
}

type game struct {
	sm   *soundManager
	sprm *spriteManager
	sf   *starfield
	time []float32 // gsls uniform [1]float32
	ship *ship

	// rocks     []*Rock
	rocks     List[*Rock]
	missiles  []*missile
	particles []particle

	RocksQt *QuadTree[RockListEl]
	//MissilesQt *QuadTree

	sW, sH int32
	gW, gH float64

	paused        bool
	cursorEnabled bool
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

	g.missiles = make([]*missile, 0, maxMissiles)
	g.RocksQt = NewQuadTree[RockListEl](0, Rect{0, 0, w, h})
	g.initMouse()
	g.paused = false
	g.sW, g.sH = w, h
	g.gW, g.gH = float64(w), float64(h)

	g.time = make([]float32, 1)
	g.sf = newStarfield(w, h, g.time)
	debug = startWithDebugOn
	g.sm = newSoundManager(startMuted)
	g.sprm = newSpriteManager()

	//g.objects = make([]*object, 0, 20)

	g.ship = newShip(float64(w/2), float64(h/2), 1000, 1000)
	g.ship.rot = 45 - 180
	generateRocks(g, noPreferredRocks)

	tprev = time.Now().Local().UnixMicro()

	vectorFont = rl.LoadFontEx("res/Vectorb.ttf", 99, nil, 0)
	rl.GenTextureMipmaps(&vectorFont.Texture)
	rl.SetTextureFilter(vectorFont.Texture, rl.FilterBilinear)

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
	iterator := gme.rocks.Iter()
	for el, ok := iterator(); ok; el, ok = iterator() {
		el.Value.Draw()
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

	iterator := gme.rocks.Iter()
	for el, ok := iterator(); ok; el, ok = iterator() {
		// go el.Value.Move(dt)
		el.Value.Move(dt)
	}
	// wg.Done()
}
func (gme *game) moveMissiles(dt float64) {
	for i := range gme.missiles { // move missiles
		// go gme.missiles[i].Move(dt)
		gme.missiles[i].Move(dt)
	}
	// wg.Done()
}

func (gme *game) buildRocksQTree() {
	gme.RocksQt.Clear()

	iterator := gme.rocks.Iter()
	for r, ok := iterator(); ok; r, ok = iterator() {
		gme.RocksQt.Insert(RockListEl{ListEl: *r})
	}
}

// var wg sync.WaitGroup

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

	gme.drawMissiles()
	gme.drawParticles()
	gme.ship.Draw()
	gme.drawStatusBar()
	gme.debugQt()
	rl.EndDrawing()

	gme.sm.doFade() // fade out sounds if needed

	// UPDATING ---------------------------------------------------------------------
	tnow = time.Now().UnixMicro()
	elapsed := tnow - tprev
	tprev = tnow
	dt := float64(elapsed) / 16666.0

	if !gme.paused {

		gme.time[0] += 0.01 // glsl uniform for starfield shader

		if debug {
			dt = 1
		}
		gme.ship.Move(dt)

		gme.ship.chargeUp() // chargeup ship

		// wg.Add(1) // Waitgroup
		gme.moveRocks(dt)
		// wg.Add(1)
		gme.moveMissiles(dt)
		gme.buildRocksQTree()
		// wg.Wait()

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

		// go gme.constrainRocks()
		// go gme.constrainMissiles()
		// go gme.animateParticles()

		gme.constrainRocks()
		gme.constrainMissiles()
		gme.animateParticles()
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
