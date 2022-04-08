package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var missileShape = []V2{{-1.25, 12.}, {-1.2, 3.12}, {-4.4, 0}, {4.4, 0}, {1.25, 3.12}, {1.25, 12.5}, {0, 20}}

type weapon struct {
	name     string
	maxCap   int
	curCap   int
	lowLimit int
	scoreMult float64
	cost     float64
}

const (
	missileNormal = iota
	missileTriple
	missileGuided
)

type missile interface {
	Draw()
	Move(*game, float64)
	getData() *aMissile
}

type aMissile struct {
	*shape
	motion

	launchSpeed float64
}

func (m *aMissile) Move(dt float64) {
	m.motion.Move(dt)
}

type normalMissile struct {
	aMissile
}

func (m *normalMissile) getData() *aMissile       { return &m.aMissile }
func (m *normalMissile) Draw()                    { m.shape.Draw(m.motion.pos, m.motion.rot, rl.DarkGray, rl.DarkGray) }
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
		_lineThick(m.pos.Add(m.speed.MulA(6)), m.targetRock.pos, 10, rl.Color{0, 100, 100, 90})
	}
	m.shape.Draw(m.motion.pos, m.motion.rot, rl.Gray, rl.Gray)
}
func (m *guidedMissile) Move(g *game, dt float64) {

	if m.life > 20 { // starts targeting after a while
		m.targetRock = nil
		var mindist = float64(2000)
		for r := g.rocks.Front(); r != nil; r = r.Next() {
			rock := r.Value.(*Rock)
			dist := m.pos.Sub(rock.pos).Len()
			v1 := rock.pos.Sub(m.pos).Norm()
			v2 := m.speed.Norm()
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
			v2 := m.speed.Norm()
			v1 := m.targetRock.pos.Sub(m.pos).Norm()
			angle := math.Atan2(v1.y*v2.x-v1.x*v2.y, v1.x*v2.x+v1.y*v2.y)
			angledeg := angle * rl.Rad2deg
			m.rotSpeed = angledeg*t/30 +
				_noise1D(uint8(m.life/2+int(m.randoffs)))*5 - 2.5 // random disturbance
			m.rot += m.rotSpeed * dt
		}
	}
	spd := m.launchSpeed
	m.speed = cs(m.rot).MulA(spd)

	speed := m.speed.MulA(dt)
	m.pos.Incr(speed)
	m.life++
}
func newMissile(pos V2, spd, rot float64) *aMissile {

	m := new(aMissile)

	m.launchSpeed = spd

	m.shape = newShape(missileShape)
	m.pos = pos
	m.rot = rot
	dir := cs(rot)
	m.speed = dir.MulA(m.launchSpeed)

	return m
}
func launchMissile(g *game, mtype int) {

	sSpd := V2len(g.ship.speed)

	switch mtype {
	case missileNormal:
		am := newMissile(g.ship.pos, sSpd+2, g.ship.rot)
		nnm := &normalMissile{aMissile: *am}
		g.missiles = append(g.missiles, nnm)
	case missileTriple:
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.pos, sSpd+1.9+rnd()*0.2, g.ship.rot-4-rnd())})
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.pos, sSpd+1.9+rnd()*0.2, g.ship.rot)})
		g.missiles = append(g.missiles, &normalMissile{aMissile: *newMissile(g.ship.pos, sSpd+1.9+rnd()*0.2, g.ship.rot+4+rnd())})
	case missileGuided:
		am := newMissile(g.ship.pos, sSpd+3, g.ship.rot)
		ngm := &guidedMissile{aMissile: *am}
		ngm.randoffs = uint8(rand.Intn(255))
		g.missiles = append(g.missiles, ngm)
	}
}
