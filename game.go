package main

import (
	"container/list"
	"fmt"
	"math"
	"sync"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	qt "rlbb/lib/quadtree"
	sm "rlbb/lib/soundmanager"
)

const (
	caption            = "test boom boom game"
	rSpeedMax          = 1
	noPreferredRocks   = 40
	PrefferredRockSize = 80
	maxRocks           = 100
	maxMissiles        = 50
	maxParticles       = 150
	FPS                = 60
	shieldsLowLimit    = 25
	forceFieldRadius   = 200
	gmMaxSensingRange  = 500 // guided missile sensing gmMaxSensingRange
	gmHalfSensingFov   = 22  // fuided missile field of view
	startWithDebugOn   = true
	startMuted         = false
)

var debug bool
var showgui bool

type game struct {
	sm   *sm.SoundManager
	sprm *spriteManager
	sf   *starfield
	time []float32 // gsls uniform [1]float32

	ship      *ship
	rocks     list.List
	missiles  []missile
	particles []particle
	RocksQt   *qt.QuadTree[*Rock]

	sW, sH        int32
	gW, gH        float64
	VisibleLights *Lighting
	paused        bool

	curWeapon int
	weapons   map[int]weapon

	vectorFont rl.Font
	pp         *PostProcess

	// global time counters
	tnow, tprev int64

	tickTock uint8
}

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

type PostProcess struct {
	shader      rl.Shader
	gamma       []float32
	iTime       []float32
	target      rl.RenderTexture2D
	iResolution []float32
	malfunct    []float32
}

func newPostprocess(w, h int32) *PostProcess {
	pp := &PostProcess{}
	pp.gamma = make([]float32, 1)
	pp.iTime = make([]float32, 1)
	pp.malfunct = make([]float32, 1)
	pp.iResolution = make([]float32, 2)

	pp.shader = rl.LoadShader("shaders/base.vs", "shaders/postprocess.fs")

	pp.iResolution[0], pp.iResolution[1] = float32(w), float32(h)
	rl.SetShaderValue(pp.shader, rl.GetShaderLocation(pp.shader, "iResolution"), pp.iResolution, rl.ShaderUniformVec2)

	pp.target = rl.LoadRenderTexture(w, h)

	return pp
}
func (pp *PostProcess) SetShaderValues() {
	pp.iTime[0] = pp.iTime[0] + 1

	rl.SetShaderValue(pp.shader, rl.GetShaderLocation(pp.shader, "gamma"), pp.gamma, rl.ShaderUniformFloat)
	rl.SetShaderValue(pp.shader, rl.GetShaderLocation(pp.shader, "iTime"), pp.iTime, rl.ShaderUniformFloat)
	rl.SetShaderValue(pp.shader, rl.GetShaderLocation(pp.shader, "glitch"), pp.malfunct, rl.ShaderUniformInt)
}
func (pp *PostProcess) Finalize() {
	rl.UnloadShader(pp.shader)
	rl.UnloadRenderTexture(pp.target)
	rl.UnloadShader(slshader)
}

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
	g.SetScreenSize(w, h)
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagVsyncHint | rl.FlagWindowMaximized)
	rl.InitWindow(w, h, caption)
	g.sm = sm.NewSoundManager(startMuted, soundFiles)

	g.sm.Msgs = map[int]sm.VoiceMsg{sAmmoLow: {Repeat: 15}, sShieldsLow: {Repeat: 17}}
	g.sm.EnableLoops(sSpace, sScore)

	g.VisibleLights = newLighting()
	g.VisibleLights.AddLight(OmniLight{V2{1440, 400}, _ColorfromRlColor(rl.Purple), 900})

	redsun = &OmniLight{V2{-100, 100}, _ColorfromRlColor(rl.Red), 300}
	g.VisibleLights.AddLight(redsun)

	g.pp = newPostprocess(w, h)

	rl.SetTargetFPS(FPS)

	g.missiles = make([]missile, 0, maxMissiles)
	g.RocksQt = qt.NewNode[*Rock](0, qt.Rect{X: 0, Y: 0, W: g.sW, H: g.sH})
	initMouse()
	g.paused = false

	g.time = make([]float32, 1)
	g.sf = newStarfield(w, h, g.time)
	debug = startWithDebugOn

	g.sprm = newSpriteManager()

	g.ship = newShip(float64(w/2), float64(h/2), 1000, 1000)
	g.ship.rot = 45 - 180
	g.VisibleLights.AddLight(g.ship.light)
	g.VisibleLights.AddLight(g.ship.slight)
	g.generateRocks(noPreferredRocks)

	g.tprev = time.Now().Local().UnixMicro()

	g.vectorFont = rl.LoadFontEx("res/Vectorb.ttf", 99, nil, 0)
	rl.GenTextureMipmaps(&g.vectorFont.Texture)
	rl.SetTextureFilter(g.vectorFont.Texture, rl.FilterBilinear)

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

func (g *game) playCyclicMessages() {
	wpn := g.weapons[g.curWeapon]

	if wpn.curCap < wpn.lowLimit {
		g.sm.PlayCyclic(sAmmoLow)
	} else {
		g.sm.StopCyclic(sAmmoLow)
	}
	if g.ship.shields < shieldsLowLimit {
		g.sm.PlayCyclic(sShieldsLow)
	} else {
		g.sm.StopCyclic(sShieldsLow)
	}
}

func (g *game) drawStatusBar() {
	if !g.paused {
		wpn := g.weapons[g.curWeapon]
		rl.DrawRectangle(0, g.sH-20, g.sW, 26, rl.DarkPurple)
		_multicolorText(20, g.sH-20, 20,
			"Cash:", rl.Purple, g.ship.cash, rl.Purple, 30, 10,
			"Shields:", rl.Purple, int(g.ship.shields), _flashColor(Game.tickTock, rl.Purple, 50, shieldsLowLimit, int(g.ship.shields)),
			"Energy:", rl.Purple, int(g.ship.energy), _flashColor(Game.tickTock, rl.Purple, 500, 100, int(g.ship.energy)),
			wpn.name+"s :", rl.Purple, wpn.curCap, _flashColor(Game.tickTock, rl.Purple, 30, wpn.lowLimit, wpn.curCap))

		rl.DrawFPS(g.sW-80, g.sH-20)
	} else {
		var col rl.Color
		if Game.tickTock%20 > 10 {
			col = rl.DarkPurple
		} else {
			col = rl.Purple
		}
		rl.DrawText("**** GAME PAUSED ***", 20, g.sH-20, 20, col)
	}
}

func (g *game) drawGUI() {
	if showgui {
		x := g.gW - 140
		y := float32(20)

		clicked := gui.Button(rl.Rectangle{float32(x), y, 140, 20}, "Clicken mich")
		str := fmt.Sprintf("das buton war cliked : %v", clicked)
		y += 21
		gui.Label(rl.Rectangle{float32(x), y, 140, 20}, str)
		y += 21

		malFunXT := g.pp.malfunct[0] != 0
		if gui.CheckBox(rl.Rectangle{float32(x) + 100, y, 20, 20}, malFunXT) {
			g.pp.malfunct[0] = 1
		} else {
			g.pp.malfunct[0] = 0
		}

		gui.Label(rl.Rectangle{float32(x), y, 100, 20}, "Malfunxon")
		y += 21
		g.pp.gamma[0] = gui.Slider(rl.Rectangle{float32(x), y, 140, 20}, g.pp.gamma[0], 0.0, 3.0)
		y += 21
		str = fmt.Sprintf("gamma:%v", g.pp.gamma[0])
		gui.Label(rl.Rectangle{float32(x), y, 100, 20}, str)
		y += 21
		gui.Label(rl.Rectangle{float32(x), y, 100, 20}, "Shields")
		y += 21
		g.ship.shields = float64(gui.Slider(rl.Rectangle{float32(x), y, 140, 20}, float32(g.ship.shields), 0.0, 100.0))
	}
}

func (gme *game) addParticle(p particle) bool {
	if len(gme.particles) < maxParticles {
		gme.particles = append(gme.particles, p)
		return true
	} else {
		return false
	}
}

func (gme *game) animateParticles() {

	var wg sync.WaitGroup
	for i := 0; i < len(gme.particles); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			gme.particles[idx].Animate()
		}(i)
	}
	wg.Wait()
	i := 0
	for i < len(gme.particles) {
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

}
func (gme *game) moveMissiles(dt float64) {
	for i := range gme.missiles { // move missiles
		// go gme.missiles[i].Move(dt)
		gme.missiles[i].Move(gme, dt)
	}

}

func (gme *game) buildRocksQTree() {
	gme.RocksQt = qt.NewNode[*Rock](0, qt.Rect{X: 0, Y: 0, W: gme.sW, H: gme.sH})

	for el := gme.rocks.Front(); el != nil; el = el.Next() {
		gme.RocksQt.Insert(el.Value.(*Rock))
	}
}

const starsNo = 1000

func (gme *game) GameDraw() {
	rl.BeginDrawing()

	rl.BeginTextureMode(gme.pp.target)

	gme.sf.draw() // draw starfield
	gme.drawForceField()
	redsun.SetColor(Color{_noise1D(Game.tickTock)*0.4 + 0.6, 0, 0, 1})

	rl.BeginBlendMode(rl.BlendAdditive)
	gme.VisibleLights.Draw()
	rl.EndBlendMode()
	gme.drawMissiles()
	gme.drawAndDeleteRocks()

	gme.drawParticles()
	
	gme.ship.Draw()

	gme.debugQt()

	rl.EndTextureMode()

	rl.BeginShaderMode(gme.pp.shader)
	gme.pp.SetShaderValues()
	rl.DrawTextureRec(gme.pp.target.Texture,
		rl.NewRectangle(0, 0, float32(gme.pp.target.Texture.Width), float32(-gme.pp.target.Texture.Height)),
		rl.NewVector2(0, 0), rl.White)
	rl.EndShaderMode()

	gme.drawStatusBar()
	gme.drawGUI()
	rl.EndDrawing()

	gme.tickTock++
}

var wg sync.WaitGroup

func (gme *game) GameUpdate(dt float64) {

	// UPDATING ---------------------------------------------------------------------
	// postprocess
	// gamma - gets darker perioddically
	// malfunct - the screen glitches depending on ship shields level (damage)
	gme.pp.gamma[0] = float32(1.0 + math.Sin(float64(gme.time[0]/30))*0.3)
	if gme.ship.shields < 80 && !gme.ship.destroyed {
		pwm := 1 - gme.ship.shields/80
		scale := gme.ship.shields / 8
		if rnd() < pwm/scale {
			gme.pp.malfunct[0] = 1.0
		} else {
			gme.pp.malfunct[0] = 0.0
		}
	} else {
		gme.pp.malfunct[0] = 0.0
	}

	if !gme.paused {

		gme.sm.Update()

		gme.time[0] += 0.01 // glsl uniform for starfield shader

		if debug {
			dt = 1
		}
		gme.ship.Move(dt)

		gme.ship.ChargeUp() // chargeup ship energy

		wg.Add(1)
		go func() {
			defer wg.Done()
			gme.moveRocks(dt)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			gme.moveMissiles(dt)
		}()

		wg.Wait() // wait for the above two procedures to complete

		gme.buildRocksQTree()

		gme.processMissileHits()
		gme.processForceField()
		gme.processShipHits()

		go gme.constrainShip()
		go gme.constrainRocks()
		go gme.constrainMissiles()
		gme.animateParticles()
	}
}

var cursorEnabled bool
var mouseLastmoved int64

func initMouse() {
	rl.DisableCursor()
	cursorEnabled = false
	mouseLastmoved = time.Now().Unix()
}
func processMouse() {
	dx, dy := rl.GetMouseDelta().X, rl.GetMouseDelta().X

	if !cursorEnabled && dx*dx+dy*dy > 16 {
		rl.EnableCursor()
		cursorEnabled = true
		mouseLastmoved = time.Now().Unix()
	} else if cursorEnabled && time.Now().Unix()-mouseLastmoved > 3 {
		rl.DisableCursor()
		cursorEnabled = false
	}
}
func (g *game) finalize() {
	rl.CloseWindow()
	g.sm.UnloadAll()
	g.sprm.unloadAll()
	g.pp.Finalize()
}
