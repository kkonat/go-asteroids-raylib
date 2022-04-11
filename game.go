package main

import (
	"container/list"
	"image/color"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	// "github.com/golang-ui/nuklear/nk"	// interface, maybe
	qt "rlbb/lib/quadtree"
	sm "rlbb/lib/soundmanager"

	"golang.org/x/exp/constraints"
)

const (
	caption            = "test boom boom game"
	rSpeedMax          = 1
	noPreferredRocks   = 40
	PrefferredRockSize = 80
	maxRocks           = 100
	maxMissiles        = 50
	maxParticles       = 50
	FPS                = 60
	shieldsLowLimit    = 25
	forceFieldRadius   = 200
	gmMaxSensingRange  = 500 // guided missile sensing gmMaxSensingRange
	gmHalfSensingFov   = 22  // fuided missile field of view
	startWithDebugOn   = true
	startMuted         = false
)

type game struct {
	sm        *sm.SoundManager
	sprm      *spriteManager
	sf        *starfield
	time      []float32 // gsls uniform [1]float32
	ship      *ship
	rocks     list.List
	missiles  []missile
	particles []particle
	RocksQt   *qt.QuadTree[*Rock]

	sW, sH        int32
	gW, gH        float64
	Lights        *Lighting
	paused        bool
	cursorEnabled bool
	curWeapon     int
	weapons       map[int]weapon
}

var vectorFont rl.Font
var debug bool

// global time counters
var tnow, tprev, mAmmoLastPlayed, mShieldsLastPlayed int64
var tickTock uint8

const (
	sSpace = iota
	sScore
	sMissilesDlvrd
	sThrust
	sExpl
	sLaunch
	sShieldsLow
	sAmmoLow
	sOinx
	sExplodeShip
	sScratch
	sChargeUp
	sForceField
)

var redsun *OmniLight

func newGame(w, h int32) *game {

	g := new(game)
	soundFiles := map[int]sm.SoundFile{
		// Id			  filename        vol  pitch
		sSpace:         {Fname: "res/space.ogg", Vol: 0.5, Pitch: 1.0},
		sScore:         {Fname: "res/score.mp3", Vol: 0.1, Pitch: 1.0},
		sMissilesDlvrd: {Fname: "res/missiles-delivered.ogg", Vol: 0.5, Pitch: 1.0},
		sThrust:        {Fname: "res/thrust.ogg", Vol: 0.5, Pitch: 1.0},
		sExpl:          {Fname: "res/expl.ogg", Vol: 0.5, Pitch: 0.65},
		sLaunch:        {Fname: "res/launch.ogg", Vol: 0.5, Pitch: 1.0},
		sShieldsLow:    {Fname: "res/warning-shields-low.ogg", Vol: 0.3, Pitch: 1.0},
		sAmmoLow:       {Fname: "res/warning-ammo-low.ogg", Vol: 0.3, Pitch: 1.0},
		sOinx:          {Fname: "res/oinxL.ogg", Vol: 0.5, Pitch: 1.0},
		sExplodeShip:   {Fname: "res/shipexplode.ogg", Vol: 1.0, Pitch: 1.0},
		sScratch:       {Fname: "res/metalScratch.ogg", Vol: 0.2, Pitch: 1.0},
		sChargeUp:      {Fname: "res/chargeup.ogg", Vol: 0.2, Pitch: 1.0},
		sForceField:    {Fname: "res/forcefield2.ogg", Vol: 0.5, Pitch: 1.0},
	}

	g.sm = sm.NewSoundManager(startMuted, soundFiles)

	g.sm.EnableLoops(sSpace, sScore)
	g.SetScreenSize(w, h)
	g.Lights = &Lighting{}
	g.Lights.AddLight(OmniLight{V2{1440, 400}, _ColorfromRlColor(rl.Purple), 900})

	redsun = &OmniLight{V2{-100, 100}, _ColorfromRlColor(rl.Red), 300}
	g.Lights.AddLight(redsun)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)
	rl.InitWindow(w, h, caption)

	rl.SetTargetFPS(FPS)

	g.missiles = make([]missile, 0, maxMissiles)
	g.RocksQt = qt.NewNode[*Rock](0, qt.Rect{X: 0, Y: 0, W: g.sW, H: g.sH})
	g.initMouse()
	g.paused = false

	g.time = make([]float32, 1)
	g.sf = newStarfield(g.sW, g.sH, g.time)
	debug = startWithDebugOn

	g.sprm = newSpriteManager()

	g.ship = newShip(float64(w/2), float64(h/2), 1000, 1000)
	g.ship.rot = 45 - 180
	g.generateRocks(noPreferredRocks)

	tprev = time.Now().Local().UnixMicro()

	vectorFont = rl.LoadFontEx("res/Vectorb.ttf", 99, nil, 0)
	rl.GenTextureMipmaps(&vectorFont.Texture)
	rl.SetTextureFilter(vectorFont.Texture, rl.FilterBilinear)

	g.weapons = make(map[int]weapon)
	g.weapons[missileNormal] = weapon{"missile", 100, 100, 20, 4.0, 1.6}
	g.weapons[missileTriple] = weapon{"triple", 100, 100, 20, 1.3, 4.8}
	g.weapons[missileGuided] = weapon{"guided missile", 20, 20, 4, 3.0, 3.2}

	return g
}
func (g *game) SetScreenSize(w, h int32) {
	g.sW, g.sH = w, h
	g.gW, g.gH = float64(w), float64(h)
}

func (g *game) playMessages() {
	wpn := g.weapons[g.curWeapon]
	if wpn.curCap < wpn.lowLimit {
		t := time.Now().Local().Unix()
		if mAmmoLastPlayed == 0 || t-mAmmoLastPlayed > 15 {
			{
				g.sm.PlayM(sAmmoLow)
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
				g.sm.PlayM(sShieldsLow)
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
		wpn := g.weapons[g.curWeapon]
		rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
		_multicolorText(20, g.sH-20, 20,
			"Cash:", rl.Purple, g.ship.cash, rl.Purple, 30, 10,
			"Shields:", rl.Purple, int(g.ship.shields), flashColor(rl.Purple, 50, shieldsLowLimit, int(g.ship.shields)),
			"Energy:", rl.Purple, int(g.ship.energy), flashColor(rl.Purple, 500, 100, int(g.ship.energy)),
			wpn.name+"s :", rl.Purple, wpn.curCap, flashColor(rl.Purple, 30, wpn.lowLimit, wpn.curCap))

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
func (gme *game) drawAndDeleteRocks() {
	for el := gme.rocks.Front(); el != nil; el = el.Next() {
		if el.Value.(*Rock).delete {
			gme.rocks.Remove(el)
		} else {
			el.Value.(*Rock).Draw()
		}
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

	for el := gme.rocks.Front(); el != nil; el = el.Next() {
		el.Value.(*Rock).Move(dt)
	}
	//	wg.Done()
}
func (gme *game) moveMissiles(dt float64) {
	for i := range gme.missiles { // move missiles
		// go gme.missiles[i].Move(dt)
		gme.missiles[i].Move(gme, dt)
	}
	//	wg.Done()
}

func (gme *game) buildRocksQTree() {
	gme.RocksQt = qt.NewNode[*Rock](0, qt.Rect{X: 0, Y: 0, W: gme.sW, H: gme.sH})

	for el := gme.rocks.Front(); el != nil; el = el.Next() {
		gme.RocksQt.Insert(el.Value.(*Rock))
	}
}

var wg sync.WaitGroup

func (gme *game) GameDraw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.Black)
	gme.sf.draw() // draw starfield
	gme.drawForceField()
	redsun.SetColor(Color{_noise1D(tickTock)*0.4 + 0.6, 0, 0, 1})
	gme.Lights.Draw()
	gme.drawMissiles()
	gme.drawAndDeleteRocks()

	gme.drawParticles()
	gme.ship.Draw()
	gme.drawStatusBar()
	gme.debugQt()
	rl.EndDrawing()

	tickTock++
}

func (gme *game) GameUpdate() {

	// UPDATING ---------------------------------------------------------------------
	tnow = time.Now().UnixMicro()
	elapsed := tnow - tprev
	tprev = tnow
	dt := float64(elapsed) / 16666.0

	if !gme.paused {

		gme.sm.Update()

		gme.time[0] += 0.01 // glsl uniform for starfield shader

		if debug {
			dt = 1
		}
		gme.ship.Move(dt)

		gme.ship.chargeUp() // chargeup ship energy

		//		wg.Add(1)
		gme.moveRocks(dt)
		//		wg.Add(1)
		gme.moveMissiles(dt)

		//wg.Wait() // wait for the above two procedures to complete

		gme.buildRocksQTree()

		gme.processMissileHits()
		gme.processForceField()
		gme.processShipHits()

		gme.constrainShip()

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
	g.sm.UnloadAll()
	g.sprm.unloadAll()
}
