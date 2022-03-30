package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type particle interface {
	Draw()
	Animate()
	canDelete() bool
}

type sparks struct {
	timer, timerMax        int
	positions, speeds      []V2int
	lives, maxlives, seeds []uint8
	life                   int
	sparksNo               int
	sCol, eCol             rl.Color
}

func newSparks(pos, mspeed V2, count int, maxradius, duration float64, sCol, eCol rl.Color) *sparks {
	s := new(sparks)
	s.sparksNo = count + rand.Intn(count/2)
	speed := 0.5 + rnd()*1.5
	s.positions = make([]V2int, s.sparksNo)
	s.speeds = make([]V2int, s.sparksNo)
	s.lives = make([]uint8, s.sparksNo)
	s.maxlives = make([]uint8, s.sparksNo)
	s.seeds = make([]uint8, s.sparksNo)
	s.sCol, s.eCol = sCol, eCol
	angle := 0.0
	frames := int(duration * FPS)

	s.life = frames
	for i := 0; i < s.sparksNo; i++ {
		angle += (360 / float64(s.sparksNo)) + rndSym(15)
		s.positions[i] = pos.ToV2int()
		sp := mspeed.Add(rotV(angle).MulA(5 * speed * (0.5 + rnd())))
		s.speeds[i] = sp.ToV2int()
		s.maxlives[i] = uint8(frames/2 + rand.Intn(frames/2))
		s.seeds[i] = uint8(rand.Intn(256))
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
		disturb := _noise2D(s.lives[i] + s.seeds[i]).MulA(age).SubA(age / 2).ToV2int()

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
func (s *sparks) Draw() {
	for i := 0; i < s.sparksNo; i++ {
		if s.lives[i] < s.maxlives[i] {
			if s.lives[i] < s.maxlives[i]/3 {
				c := _colorBlend(s.lives[i], s.maxlives[i]/3, s.eCol, s.sCol)
				_rectV2int(s.positions[i], 2, c)
			} else {
				t := float32(s.lives[i]-s.maxlives[i]/3) / (float32(s.maxlives[i] / 3 * 2))
				v := float32(rand.Intn(2))
				c := rl.ColorFromHSV(t*33, 1.0, v)

				_rectV2int(s.positions[i], 2, c) // I assume this is  faster than circle
			}
		}
	}
}

type explosion struct {
	timer, timerMax       int
	position, speed, offs V2
	r, rstep              float64
	maxr, dur, t          float64
}

func newExplosion(pos, speed V2, maxradius, duration float64) *explosion {
	e := new(explosion)
	e.position = pos
	e.speed = speed
	e.offs = e.position.Add(V2{rndSym(maxradius / 10), rndSym(maxradius / 10)})
	e.maxr, e.dur = maxradius, duration
	e.rstep = maxradius / (duration * FPS)
	e.timerMax = int(duration * FPS)
	return e
}

func (e *explosion) Animate() {
	if e.timer < e.timerMax {
		e.timer++
	}
	e.position = e.position.Add(e.speed)
	e.offs = e.offs.Add(e.speed)
	e.r += e.rstep
}

func (e *explosion) canDelete() bool {
	return e.timer >= e.timerMax
}

func (e *explosion) Draw() {
	t := 1 - e.timer/(e.timerMax/2)
	if e.timer < e.timerMax/3 { // phase 1 - flash
		_gradientdisc(e.position, e.r*e.r*e.r/5, rl.ColorAlpha(rl.Yellow, float32(t)*0.3), rl.Black)
	}
	if e.timer < e.timerMax/2 { // phase 2 - fireball
		_disc(e.position, e.r, rl.Yellow)

	} else { // phase 3 black grow
		r := e.r*2 - e.r/2
		t := float32((e.timer - e.timerMax/2)) / float32(e.timerMax/2)
		_gradientdisc(e.position, e.maxr, rl.ColorAlpha(rl.Yellow, 1-t), rl.ColorAlpha(rl.Orange, 1-t))
		_disc(e.offs, r, rl.Black)
	}
}
