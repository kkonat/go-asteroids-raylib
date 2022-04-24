package main

import (
	"log"
	"math/rand"
	v "rlbb/lib/vector"

	"github.com/fogleman/ease"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type particle interface {
	Draw()
	Animate()
	canDelete() bool
}

type textPart struct {
	text          string
	textW         int32
	size, growSpd float32
	pos, speed    V2
	life, maxLife uint8
	sCol, eCol    rl.Color
}

// emits new text particle with given pos, speed, text, duration, growSpeed, randomDir, start and end color
func newTextPart(pos, speed V2, text string, size int32, duration, growSpd float32, randomDir bool, sCol, eCol rl.Color) *textPart {
	tp := new(textPart)
	tp.speed = speed.MulA(0.5)
	tp.growSpd = growSpd
	tp.pos = pos
	tp.text = text
	if randomDir {
		tp.speed = v.RotV(rnd() * 360)
	}
	tp.size = float32(size)
	tp.life = uint8(duration * FPS)
	tp.maxLife = tp.life
	tp.sCol = sCol
	tp.eCol = rl.Fade(eCol, 0)
	return tp
}
func (s *textPart) canDelete() bool {
	if s.life > 0 {
		return false
	} else {
		return true
	}
}
func (tp *textPart) Animate() {
	tp.pos.Incr(tp.speed)
	tp.speed.MulA(0.95)
	tp.size += tp.growSpd
	if tp.life > 0 {
		tp.life--
	}
}

// draws text particle
func (tp *textPart) Draw() {
	col := _colorBlend(tp.life, tp.maxLife, tp.eCol, tp.sCol) // tp.life goes from 1 to 0, so reverse blend
	textw := rl.MeasureText(tp.text, int32(tp.size)) / 2
	rl.DrawText(tp.text, int32(tp.pos.X)-textw, int32(tp.pos.Y-float64(tp.size)/2), int32(tp.size), col)

}

type sparks struct {
	timer, timerMax        int
	positions, speeds      []v.FxdFPV2
	lives, maxlives, seeds []uint8
	life                   int
	sparksNo               int
	sCol, eCol             rl.Color
}

// emits  sparks particles
func newSparks(pos, mspeed V2, speedmul float64, count int, maxradius, duration float64, sCol, eCol rl.Color) *sparks {
	s := new(sparks)
	s.sparksNo = count + rand.Intn(count/2)
	speed := 0.5 + rnd()*1.5
	s.positions = make([]v.FxdFPV2, s.sparksNo)
	s.speeds = make([]v.FxdFPV2, s.sparksNo)
	s.lives = make([]uint8, s.sparksNo)
	s.maxlives = make([]uint8, s.sparksNo)
	s.seeds = make([]uint8, s.sparksNo)
	s.sCol, s.eCol = sCol, eCol
	angle := 0.0
	frames := int(duration * FPS)

	s.life = frames
	for i := 0; i < s.sparksNo; i++ {
		angle += (360 / float64(s.sparksNo)) + rndSym(15)
		s.positions[i] = pos.ToFxdFPV2()
		sp := mspeed.Add(v.RotV(angle).MulA(5 * speed * speedmul * rnd()))
		s.speeds[i] = sp.ToFxdFPV2()
		s.positions[i] = s.positions[i].Add(s.speeds[i])
		s.maxlives[i] = uint8(frames/2 + rand.Intn(frames/2))
		s.seeds[i] = uint8(rand.Intn(256))
		if s.seeds[i] == 0 {
			s.seeds[i] = 1
		}
	}

	return s
}

func (s *sparks) canDelete() bool {
	if s.life > 0 {
		return false
	} else {
		return true
	}
}

const damp = 255 // (1<<8) * 0.996

func (s *sparks) Animate() {
	for i := 0; i < s.sparksNo; i++ {
		age := float64(s.lives[i]) / 10
		disturb := _noise2D(s.lives[i] + s.seeds[i]).MulA(age).SubA(age / 2).ToFxdFPV2()

		s.speeds[i] = s.speeds[i].MulA(damp)
		s.positions[i] = s.positions[i].Add(s.speeds[i].Add(disturb))
		if s.lives[i] < s.maxlives[i] {
			s.lives[i]++
		}
	}
	if s.life > 0 {
		s.life--

	}
}

var sprkshader rl.Shader
var sprkimg *rl.Image
var sprktextr rl.Texture2D

func initParticleShaders() {
	sprkshader = rl.LoadShader("shaders/base.vs", "shaders/spark.fs")
	sprkimg = rl.GenImageColor(64, 64, rl.Black)
	sprktextr = rl.LoadTextureFromImage(sprkimg)
	//rl.SetShaderValue(sprkshader, rl.GetShaderLocation(sprkshader, "iResolution"), iResolution, rl.ShaderUniformVec2)
}
func deinitParticleShareds() {
	rl.UnloadShader(sprkshader)
	rl.UnloadTexture(sprktextr)
	rl.UnloadImage(sprkimg)

	
}

func (s *sparks) Draw() {
	rl.BeginShaderMode(sprkshader)
	rl.BeginBlendMode(rl.BlendAdditive)
	for i := 0; i < s.sparksNo; i++ {
		if s.lives[i] < s.maxlives[i] {
			if s.lives[i] < s.maxlives[i]/3 {
				c := _colorBlend(s.lives[i], s.maxlives[i]/3, s.sCol, s.eCol)

				rl.DrawTexture(sprktextr, s.positions[i].X.ToInt32(), s.positions[i].Y.ToInt32(), c)
				//DrawRectangleV(rl.Vector2{0, 0}, rl.Vector2{float32(Game.gW), float32(Game.gH)}, c)
				//_rectFxdV2(s.positions[i], 2, c)
			} else {
				t := float32(s.lives[i]-s.maxlives[i]/3) / (float32(s.maxlives[i] / 3 * 2))
				//v := float32(rand.Intn(2))
				c := _colorBlendFloat(t, s.eCol, rl.Fade(s.eCol, t))
				//c := rl.ColorFromHSV(t*33, 1.0, v)

				rl.DrawTexture(sprktextr, s.positions[i].X.ToInt32(), s.positions[i].Y.ToInt32(), c)
				//_rectFxdV2(s.positions[i], 2, c)
			}
		}
	}
	rl.EndBlendMode()
	rl.EndShaderMode()
}

type explosion struct {
	timer, timerMax       int
	position, speed, offs V2
	r, rstep              float64
	maxr, dur, t          float64
	light                 *OmniLight
}

const explosionLightStrength = 250

// emits new missile explosion particle
func newExplosion(pos, speed V2, maxradius, duration float64) *explosion {
	e := new(explosion)
	e.position = pos
	e.speed = speed
	e.offs = e.position.Add(V2{X: rndSym(maxradius / 10), Y: rndSym(maxradius / 10)})
	e.maxr, e.dur = maxradius, duration
	e.rstep = maxradius / (duration * FPS)
	e.timerMax = int(duration * FPS)
	e.light = &OmniLight{pos, Color{0.78, 0.78, 0, 1.0}, explosionLightStrength}

	return e
}

func (e *explosion) Animate() {
	if e.timer < e.timerMax {
		e.timer++
	}
	e.position = e.position.Add(e.speed)
	if e.light != nil {
		e.light.Pos = e.position
	}
	e.offs = e.offs.Add(e.speed)
	e.r += e.rstep
}

func (e *explosion) canDelete() bool {
	return e.timer >= e.timerMax
}

func (e *explosion) Draw() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	t := 1 - float64(e.timer)/float64(e.timerMax/2)
	if e.timer < e.timerMax/3 { // phase 1 - flash
		_gradientdisc(e.position, e.r*e.r*e.r/5, rl.ColorAlpha(rl.Yellow, float32(t)*0.3), rl.Black)

	}
	if e.timer < e.timerMax/2 { // phase 2 - fireball
		rl.BeginBlendMode(rl.BlendAdditive)
		_disc(e.position, e.r, rl.Yellow)
		rl.EndBlendMode()
		e.light.Strength = ease.OutBounce(float64(1-t)) * explosionLightStrength

	} else { // phase 3 black grow
		r := e.r*2 - e.r/2
		t := float32((e.timer - e.timerMax)) / float32(e.timerMax)
		rl.BeginBlendMode(rl.BlendAdditive)
		_gradientdisc(e.position, e.maxr, rl.ColorAlpha(rl.Yellow, 1-t), rl.ColorAlpha(rl.Orange, 1-t))
		rl.EndBlendMode()
		_gradientdisc(e.offs, r, rl.Fade(rl.Black, 1.0), rl.Fade(rl.Black, float32(ease.OutQuint(float64(1.0-t)))))
		//_disc(e.offs, r, rl.Black)
		strength := explosionLightStrength * (1 - e.r*e.r/(r*r))
		e.light.Strength = float64(strength / 2)
	}
	if e.timer >= e.timerMax-1 {
		Game.VisibleLights.DeleteLight(e.light)
	}
}
