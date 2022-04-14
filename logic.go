package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	qt "rlbb/lib/quadtree"
	v "rlbb/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const minidist2 = PrefferredRockSize * PrefferredRockSize * 16

func (g *game) constrainShip() {
	const limit = 40.0
	const getback = 0.5
	if g.ship.pos.X < limit || g.ship.pos.X > g.gW-limit || g.ship.pos.Y < limit || g.ship.pos.Y > g.gH-limit {
		if g.ship.isSliding {
			g.ship.speed = g.ship.speed.MulA(0.9)
			if g.ship.pos.Len2() < 0.01 {
				g.ship.isSliding = false
			}
		} else {
			if g.ship.pos.X < limit {
				g.ship.speed.X = getback
			}
			if g.ship.pos.X > g.gW-limit {
				g.ship.speed.X = -getback
			}
			if g.ship.pos.Y < limit {
				g.ship.speed.Y = getback
			}
			if g.ship.pos.Y > g.gH-limit {
				g.ship.speed.Y = -getback
			}
		}

	}
}

func (g *game) constrainRocks() {
	const limit = -50 // ordinate to detect rock going off screen / respawn new one

	var count int
	for rock := g.rocks.Front(); rock != nil; rock = rock.Next() {
		count++
		r := rock.Value.(*Rock)
		if (r.pos.X+r.radius < limit && r.speed.X < 0) ||
			(r.pos.X-r.radius > g.gW-limit && r.speed.X > 0) ||
			(r.pos.Y+r.radius < limit && r.speed.Y < 0) ||
			(r.pos.Y-r.radius > g.gH-limit && r.speed.Y > 0) {
			if g.rocks.Len() < noPreferredRocks {
				// respawn rock in a new sector, first randomly on the screen
				r.randomize()
				r.speed = V2{X: rnd()*rSpeedMax - rSpeedMax/2, Y: rnd()*rSpeedMax - rSpeedMax/2}
				r.pos = V2{X: rnd() * g.gW, Y: rnd() * g.gH}

				p := &r.pos     // these variables addes to make the
				s := &r.speed   // switch statement below more redable
				rad := r.radius // kind of hack'ish but reads easier

				sect := rand.Intn(4) // random sector from wchich new rock is to originate

				switch sect { // move the rock off the screen
				case 0: // left
					{
						p.X = -rad + limit //
						if s.X < 0 {
							s.X = -s.X
						}
					}
				case 1: // top
					{
						p.Y = -rad + limit
						if s.Y < 0 {
							s.X = -s.X
						}
					}
				case 2: // right
					{
						p.X = g.gW + rad - limit
						if s.X > 0 {
							s.X = -s.X
						}
					}
				case 3: // down
					{
						p.Y = g.gH + rad - limit
						if s.Y > 0 {
							s.Y = -s.Y
						}
					}
				}
			} else {
				if rock == nil || g.rocks.Len() == 0 {
					panic("delete rock outside")
				}
				g.rocks.Remove(rock)
			}
		}
	}
	//	fmt.Println("Rock count:", count, " rocklist len=", g.rocks.Len())
}
func (g *game) constrainMissiles() {
	const limit = -10
	var i int
	for i < len(g.missiles) {
		p := g.missiles[i].getData().pos
		if p.X <= limit || p.X > g.gW-limit ||
			p.Y <= limit || p.Y > g.gH-limit {

			g.deleteMissile(i)

		} else {
			i++
		}
	}
}

func (g *game) deleteMissile(m int) {
	if m > len(g.missiles) || m < 0 { // DEBUG
		log.Print("bad missile delete ")
	} else {
		g.missiles = append(g.missiles[:m], g.missiles[m+1:]...)
	}

}

type circle struct {
	rect qt.Rect
	p    V2
	r    float64
}

func newCircleV2(p V2, r float64) *circle {
	return &circle{
		qt.Rect{X: int32(p.X - r), Y: int32(p.Y - r),
			W: int32(r * 2), H: int32(r * 2)},
		p, r}
}
func (g *game) drawForceField() {
	if g.ship.forceField {
		c := uint8(0)
		_circleGradient(g.ship.pos, forceFieldRadius,
			rl.Fade(color.RGBA{0, 200, 200, 255}, 0.25),
			rl.Fade(color.RGBA{0, 240, 200, 255}, 0.0))
		for a := float64(10) * rnd(); a < 360+float64(10)*rnd(); a += 1 + 10*rnd() {
			p := v.Cs(a).MulA(float64(forceFieldRadius) * rnd())
			_lineThick(g.ship.pos, g.ship.pos.Add(p), rand.Float32()*20+10,
				rl.Fade(color.RGBA{0, 100, 100, 127}, 0.05))
			c++
		}
	}
}
func (g *game) processForceField() {

	if g.ship.forceField {
		dSpeed := V2{}

		for r := g.rocks.Front(); r != nil; r = r.Next() {
			rock := r.Value.(*Rock)
			v := g.ship.pos.Sub(rock.pos)
			dist := v.Len() - rock.radius
			if dist < forceFieldRadius {
				dSpeed.Incr(v.Norm().MulA(2000).DivA(dist + 0.01))
			}
		}
		d := dSpeed.MulA(.001)
		ss := g.ship.speed.Len()
		d = d.DivA(ss + 0.1)
		v := d.Len()
		if v > 0.05 { // limit bounce speed
			v = v - 0.05
			d = d.Norm().MulA(v)
		}
		g.ship.speed.Incr(d)
	}
}

func (g *game) processShipHits() {
	v := v.Cs(g.ship.rot)

	var circles [3]*circle

	circles[0] = newCircleV2(g.ship.pos.Sub(v.MulA(10)), 8)
	circles[1] = newCircleV2(g.ship.pos.Add(v.MulA(2)), 6)
	circles[2] = newCircleV2(g.ship.pos.Add(v.MulA(10)), 4)

	// for _, c := range circles {
	// 	_circle(c.p, c.r, rl.Yellow)
	// }

	shipBRect := g.ship.shape.bRect

	potCols := g.RocksQt.MayCollide(shipBRect, minidist2)

	for _, r := range potCols {

		for _, c := range circles {
			dist2 := r.pos.Sub(c.p).Len2()
			if dist2 < squared(r.radius+c.r) {
				// _disc(c.p, c.r, rl.Red)
				if g.ship.shields > 0.7 {
					if !debug {
						g.ship.shields -= 0.7
					}
					g.sm.PlayFor(sScratch, 80)
				} else {
					if !debug {
						g.addParticle(newSparks(g.ship.pos, g.ship.speed, 200, 260, 5, rl.White, rl.Red))
						malFunXT = false
						if !g.ship.destroyed {

							g.sm.Play(sExplodeShip)
						}
						g.ship.destroyed = true
					}
					//game_over()
				}
			} else {
				g.sm.Stop(sScratch)
			}

		}
	}
}

var cycle, corrupted int

func (g *game) checkIfOntheList(ptcls []*Rock) (int, any) {
	var non int
	nonl := make([]*Rock, 0)
	for i := range ptcls {
		found := false
		for el := g.rocks.Front(); el != nil; el = el.Next() {
			if ptcls[i] == el.Value {
				found = true
				break
			}
		}
		if !found {
			non++
			nonl = append(nonl, ptcls[i])
		}
	}
	if non > 0 {
		return non, nonl[0]
	} else {
		return 0, nil
	}
}
func (g *game) processMissileHits() {
	const mr = 10 // missile radius
	var hit bool

	mi := 0
	for mi < len(g.missiles) {
		missile := g.missiles[mi]
		mp := missile.getData().pos
		missileBRect := missile.getData().shape.bRect

		potCols := g.RocksQt.MayCollide(missileBRect, minidist2)

		if debug {
			if degubDrawMissileLines {
				for _, c := range potCols {
					_line(mp, c.pos, rl.Red)
				}
			}
		}

		hit = false
		for i := range potCols {
			r := potCols[i]

			rp := r.pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(r.radius+mr) && !r.delete { // hit
				hitR := potCols[i]

				distBonus := g.ship.pos.Sub(r.pos).Len() / 200
				score := 1 + int((100/r.radius*distBonus/3)*g.weapons[g.curWeapon].scoreMult/4)
				g.ship.cash += score
				str := fmt.Sprintf("+%d", score)
				g.addParticle(newTextPart(missile.getData().pos, missile.getData().speed, str, 16, 2, 0, false, rl.Yellow, rl.Red))
				if expl := newExplosion(missile.getData().pos, missile.getData().speed, 30, 0.5); g.addParticle(expl) {
					Game.VisibleLights.AddLight(expl.light) // only add light
				}
				g.addParticle(newSparks(missile.getData().pos, missile.getData().speed, 100, 100, 2.0, rl.Orange, rl.Red))

				// sound
				g.sm.PlayPM(sExpl, 0.5+rnd32())

				hitR.delete = true

				// split rock
				nr := r.split(missile.getData().pos, missile.getData().speed, 6)

				// copy new rocks
				for i := 0; i < len(nr); i++ {
					if nr[i].radius > 8 {
						if g.rocks.Len() < maxRocks {
							nr[i].buildShape()
							g.rocks.PushBack(nr[i])
							g.RocksQt.Insert(nr[i])
						}
					}
				}
				g.deleteMissile(mi)
				hit = true
				break
			}
		}
		if !hit {
			mi++
		}
	}
}
