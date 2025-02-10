package main

import (
	v "bangbang/lib/vector"
	"image/color"
	"math"
	"math/rand"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var missileShape = []v.V2{
	{X: -1.25, Y: 12.}, {X: -1.2, Y: 3.12},
	{X: -4.4, Y: 0}, {X: 4.4, Y: 0},
	{X: 1.25, Y: 3.12}, {X: 1.25, Y: 12.5},
	{X: 0, Y: 20}}

type weapon struct {
	name      string
	maxCap    int
	curCap    int
	lowLimit  int
	scoreMult float64
	cost      float64
}

const (
	missileNormal = iota
	missileTriple
	missileGuided
)
const (
	flareNormal = iota
	flareSearchlight
)

type missile interface {
	Draw()
	Move(*game, float64)
	getData() *aMissile
}

type aMissile struct {
	*shape
	Motion

	launchSpeed float64
}

func (m *aMissile) Move(dt float64) {
	m.Motion.Move(dt)
}

type normalMissile struct {
	aMissile
}

func (m *normalMissile) getData() *aMissile { return &m.aMissile }
func (m *normalMissile) Draw() {

	m.shape.Draw(m.Pos, m.Rot)

	idx := uint8(int(Game.tickTock) + *((*int)(unsafe.Pointer(m))))
	disturb := _noise2D(idx * 4).MulA(6).SubA(3)
	p1 := m.Motion.Pos
	p2 := p1.Sub(m.Motion.Speed.MulA(3)).Add(disturb)
	bl := _noise1D(uint8(idx))
	c := _colorBlendA(bl, color.RGBA{30, 10, 0, 255}, color.RGBA{190, 100, 0, 255})
	_lineThick(p1, p2, 3.1, c)
}
func (m *normalMissile) Move(_ *game, dt float64) { m.aMissile.Move(dt) }

type guidedMissile struct {
	aMissile
	life       int
	randoffs   uint8
	targetRock *Rock
}

func (m *guidedMissile) getData() *aMissile { return &m.aMissile }
func (m *guidedMissile) Draw() {
	if m.targetRock != nil {
		_lineThick(m.Pos.Add(m.Speed.MulA(6)), m.targetRock.Pos, 10, rl.Color{0, 100, 100, 30})
	}
	m.shape.Draw(m.Pos, m.Rot)

	idx := uint8(m.life + *((*int)(unsafe.Pointer(m))))
	disturb := _noise2D(idx * 4).MulA(6).SubA(3)
	p1 := m.Motion.Pos
	p2 := p1.Sub(m.Motion.Speed.MulA(3)).Add(disturb)
	bl := _noise1D(uint8(idx))
	c := _colorBlendA(bl, color.RGBA{100, 50, 0, 255}, color.RGBA{190, 190, 0, 255})
	_lineThick(p1, p2, 3.1, c)
}
func (m *guidedMissile) Move(g *game, dt float64) {

	if m.life > 20 { // starts targeting after a while
		m.targetRock = nil
		var mindist = float64(2000)
		for r := g.rocks.Front(); r != nil; r = r.Next() {
			rock := r.Value.(*Rock)
			dist := m.Pos.Sub(rock.Pos).Len()
			v1 := rock.Pos.Sub(m.Pos).Norm()
			v2 := m.Speed.Norm()
			angle := math.Acos(v1.NormDot(v2))
			angularRockSize := rock.radius / dist                     // angular width of te rock
			if angle < gmHalfSensingFov*rl.Deg2rad+angularRockSize && // if missile sees it
				dist < gmMaxSensingRange && // if within missile's sensing range
				dist < mindist { // if better than previous
				mindist = dist
				m.targetRock = rock
			}

		}
		t := max(float64(m.life-20)/60, 1) // gradually blend guiding

		if m.targetRock != nil {
			v2 := m.Speed.Norm()
			v1 := m.targetRock.Pos.Sub(m.Pos).Norm()
			angle := math.Atan2(v1.Y*v2.X-v1.X*v2.Y, v1.X*v2.X+v1.Y*v2.Y)
			angledeg := angle * rl.Rad2deg
			m.RotSpeed = angledeg*t/30 +
				_noise1D(uint8(m.life/2+int(m.randoffs)))*5 - 2.5 // random disturbance
			m.Rot += m.RotSpeed * dt
		}
	}
	spd := m.launchSpeed
	m.Speed = v.RotV(m.Rot).MulA(spd)

	speed := m.Speed.MulA(dt)
	m.Pos.Incr(speed)
	m.life++
}
func newMissile(pos V2, spd, rot float64) *aMissile {

	m := new(aMissile)

	m.launchSpeed = spd

	m.shape = NewShape(missileShape, rl.Black, rl.DarkGray)
	m.Pos = pos
	m.Rot = rot
	dir := v.RotV(rot)
	m.Speed = dir.MulA(m.launchSpeed)

	return m
}
func (g *game) launchMissile() {
	mtype := g.curWeapon
	sSpd := g.ship.Speed.Len()

	switch mtype {
	case missileNormal:
		am := newMissile(g.ship.Pos, sSpd+2, g.ship.Rot)
		nnm := &normalMissile{aMissile: *am}
		g.missiles = append(g.missiles, nnm)
	case missileTriple:
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.Pos, sSpd+1.9+rnd()*0.2, g.ship.Rot-4-rnd())})
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.Pos, sSpd+1.9+rnd()*0.2, g.ship.Rot)})
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.Pos, sSpd+1.9+rnd()*0.2, g.ship.Rot+4+rnd())})
	case missileGuided:
		am := newMissile(g.ship.Pos, sSpd+3, g.ship.Rot)
		ngm := &guidedMissile{aMissile: *am}
		ngm.randoffs = uint8(rand.Intn(255))
		g.missiles = append(g.missiles, ngm)
	}
}

func (g *game) launchFlare() {}
