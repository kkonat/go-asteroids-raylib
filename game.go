package main

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	// "github.com/golang-ui/nuklear/nk"
	"golang.org/x/exp/constraints"
)

const (
	caption          = "test boom boom game"
	rSpeedMax        = 1
	noPreferredRocks = 12
	maxRocks         = 100
	maxMissiles      = 50
	maxParticles     = 50
	FPS              = 60
	shieldsLowLimit  = 25
	ammoLowLimit     = 30
)

type object interface {
	isAt() bool
	Draw()
}
type crateObj struct {
	m        motion
	spriteId int
}

func (o *crateObj) isAt() bool {
	panic("TODO")
	return false
}
func (o *crateObj) Draw() {
	panic("TODO")

}
func newCrate(spriteIdx int, pos, speed V2, rotSpeed float64) object {

	o := new(crateObj)
	o.m.pos = pos
	o.m.speed = speed
	o.spriteId = spriteIdx

	return o
}

type game struct {
	sm          *soundManager
	sprm        *spriteManager
	sf          *starfield
	ship        *ship
	rocks       [maxRocks]*Rock
	rocksNo     int
	missiles    [maxMissiles]*missile
	missilesNo  int
	particles   [maxParticles]particle
	particlesNo int
	objects     []*object

	sW, sH int32
	gW, gH float64
	ufo    rl.Texture2D
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
	g.sW, g.sH = w, h
	g.gW, g.gH = float64(w), float64(h)

	g.sf = newStarfield(w, h)

	g.sm = newSoundManager(true)
	g.sprm = newSpriteManager()

	g.objects = make([]*object, 0, 20)

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

	rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
	_multicolorText(20, g.sH-20, 20,
		"Cash:", rl.Purple, g.ship.cash, rl.Purple, 30, 10,
		"Shields:", rl.Purple, int(g.ship.shields), flashColor(rl.Purple, 50, shieldsLowLimit, int(g.ship.shields)),
		"Fuel:", rl.Purple, int(g.ship.fuel), flashColor(rl.Purple, 500, 100, int(g.ship.fuel)),
		"Missiles:", rl.Purple, g.ship.missiles, flashColor(rl.Purple, 30, ammoLowLimit, g.ship.missiles))

	rl.DrawFPS(g.sW-80, g.sH-20)

	tickTock++
}

func (gme *game) addParticle(p particle) {
	if gme.particlesNo < maxParticles {
		gme.particles[gme.particlesNo] = p
		gme.particlesNo++
	}
}
func (gme *game) animateParticles() {
	for i := 0; i < gme.particlesNo; i++ {
		gme.particles[i].Animate()
		if gme.particles[i].canDelete() {
			gme.particlesNo--

			gme.particles[i] = gme.particles[gme.particlesNo]
			gme.particles[gme.particlesNo] = nil
		}
	}
}
func (gme *game) drawRocks() {
	for i := 0; i < gme.rocksNo; i++ { // draw rocks
		gme.rocks[i].Draw()
	}
}
func (gme *game) drawMissiles() {
	for i := 0; i < gme.missilesNo; i++ { // draw missiles
		gme.missiles[i].Draw()

	}
}
func (gme *game) drawParticles() {
	for i := 0; i < gme.particlesNo; i++ {
		gme.particles[i].Draw()
	}
}

func (gme *game) moveRocks(dt float64) {
	for i := 0; i < gme.rocksNo; i++ { // move rocks
		go gme.rocks[i].m.Move(dt)
	}
	wg.Done()
}
func (gme *game) moveMissiles(dt float64) {
	for i := 0; i < gme.missilesNo; i++ { // move missiles
		go gme.missiles[i].m.Move(dt)
	}
	wg.Done()
}

var wg sync.WaitGroup

func (gme *game) drawAndUpdate() {

	// DRAWING ---------------------------------------------------------------------
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)

	gme.sf.draw() // draw starfield
	gme.drawRocks()
	gme.drawMissiles()
	gme.drawParticles()
	gme.ship.Draw() // draw ship

	gme.drawStatusBar() // draw status bar on top of everything

	rl.EndDrawing()

	gme.sm.doFade() // fade out sounds if needed

	// UPDATING ---------------------------------------------------------------------
	tnow = time.Now().UnixMicro()
	elapsed := tnow - tprev
	tprev = tnow
	dt := float64(elapsed) / 16666.0

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

func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.unloadAll()
	g.sprm.unloadAll()
}
