package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	qt "rlbb/lib/quadtree"
	v "rlbb/lib/vector"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const minidist2 = PrefferredRockSize * PrefferredRockSize * 16

func (g *game) constrainShip() {
	const limit = 40.0
	const getback = 0.5
	if g.ship.Pos.X < limit || g.ship.Pos.X > g.gW-limit || g.ship.Pos.Y < limit || g.ship.Pos.Y > g.gH-limit {
		if g.ship.isSliding {
			g.ship.Speed = g.ship.Speed.MulA(0.9)
			if g.ship.Pos.Len2() < 0.01 {
				g.ship.isSliding = false
			}
		} else {
			if g.ship.Pos.X < limit {
				g.ship.Speed.X = getback
			}
			if g.ship.Pos.X > g.gW-limit {
				g.ship.Speed.X = -getback
			}
			if g.ship.Pos.Y < limit {
				g.ship.Speed.Y = getback
			}
			if g.ship.Pos.Y > g.gH-limit {
				g.ship.Speed.Y = -getback
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
		if (r.Pos.X+r.radius < limit && r.Speed.X < 0) ||
			(r.Pos.X-r.radius > g.gW-limit && r.Speed.X > 0) ||
			(r.Pos.Y+r.radius < limit && r.Speed.Y < 0) ||
			(r.Pos.Y-r.radius > g.gH-limit && r.Speed.Y > 0) {
			if g.rocks.Len() < noPreferredRocks {
				// respawn rock in a new sector, first randomly on the screen
				r.randomize()
				r.Speed = V2{X: rnd()*rSpeedMax - rSpeedMax/2, Y: rnd()*rSpeedMax - rSpeedMax/2}
				r.Pos = V2{X: rnd() * g.gW, Y: rnd() * g.gH}

				p := &r.Pos     // these variables addes to make the
				s := &r.Speed   // switch statement below more redable
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
		p := g.missiles[i].getData().Pos
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
		_circleGradient(g.ship.Pos, forceFieldRadius,
			rl.Fade(color.RGBA{0, 200, 200, 255}, 0.25),
			rl.Fade(color.RGBA{0, 240, 200, 255}, 0.0))
		for a := float64(10) * rnd(); a < 360+float64(10)*rnd(); a += 1 + 10*rnd() {
			p := v.RotV(a).MulA(float64(forceFieldRadius) * rnd())
			_lineThick(g.ship.Pos, g.ship.Pos.Add(p), rand.Float32()*20+10,
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
			v := g.ship.Pos.Sub(rock.Pos)
			dist := v.Len() - rock.radius
			if dist < forceFieldRadius {
				dSpeed.Incr(v.Norm().MulA(2000).DivA(dist + 0.01))
			}
		}
		d := dSpeed.MulA(.001)
		ss := g.ship.Speed.Len()
		d = d.DivA(ss + 0.1)
		v := d.Len()
		if v > 0.05 { // limit bounce speed
			v = v - 0.05
			d = d.Norm().MulA(v)
		}
		g.ship.Speed.Incr(d)
	}
}

func (g *game) processShipHits() {
	v := v.RotV(g.ship.Rot)

	var circles [3]*circle

	circles[0] = newCircleV2(g.ship.Pos.Sub(v.MulA(10)), 8)
	circles[1] = newCircleV2(g.ship.Pos.Add(v.MulA(2)), 6)
	circles[2] = newCircleV2(g.ship.Pos.Add(v.MulA(10)), 4)

	// for _, c := range circles {
	// 	_circle(c.p, c.r, rl.Yellow)
	// }

	shipBRect := g.ship.shape.bRect

	potCols := g.RocksQt.MayCollide(shipBRect, minidist2)

	for _, r := range potCols {

		for _, c := range circles {
			dist2 := r.Pos.Sub(c.p).Len2()
			if dist2 < squared(r.radius+c.r) {
				// _disc(c.p, c.r, rl.Red)
				if g.ship.shields > 0.7 {
					if !debug {
						g.ship.shields -= 0.7
					}
					g.sm.PlayFor(sScratch, 80)
				} else {
					if !debug {
						if !g.ship.destroyed {
							go func() {
								for i := 0; i < 3; i++ {
									g.addParticle(newSparks(g.ship.Pos.Add(V2{rndSym(15), rndSym(15)}), g.ship.Speed, 3, 1900, 550, 4, rl.White, rl.Red))
									time.Sleep(time.Duration(111) * time.Millisecond)
								}
								time.Sleep(time.Duration(66) * time.Millisecond)
								g.addParticle(newSparks(g.ship.Pos, g.ship.Speed, 6, 10600, 990, 6, rl.White, rl.Red))
								g.sm.Play(sExplodeShip)
							}()
							g.ship.Destroy()
						}
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

func (g *game) processMissileHits() {
	const mr = 10 // missile radius
	var hit bool

	mi := 0
	for mi < len(g.missiles) {
		missile := g.missiles[mi]
		mp := missile.getData().Pos
		missileBRect := missile.getData().shape.bRect

		potCols := g.RocksQt.MayCollide(missileBRect, minidist2)

		if debug {
			if degubDrawMissileLines {
				for _, c := range potCols {
					_line(mp, c.Pos, rl.Red)
				}
			}
		}

		hit = false
		for i := range potCols {
			r := potCols[i]

			rp := r.Pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(r.radius+mr) && !r.delete { // hit
				hitR := potCols[i]

				distBonus := g.ship.Pos.Sub(r.Pos).Len() / 200
				score := 1 + int((100/r.radius*distBonus/3)*g.weapons[g.curWeapon].scoreMult/4)
				g.ship.cash += score
				str := fmt.Sprintf("+%d", score)
				g.addParticle(newTextPart(missile.getData().Pos, missile.getData().Speed, str, 16, 2, 0, false, rl.Yellow, rl.Red))
				if expl := newExplosion(missile.getData().Pos, missile.getData().Speed, 30, 0.5); g.addParticle(expl) {
					Game.VisibleLights.AddLight(expl.light) // only add light
				}
				g.addParticle(newSparks(missile.getData().Pos, missile.getData().Speed, 1, 100, 100, 2.0, rl.Orange, rl.Red))

				// sound
				g.sm.PlayPM(sExpl, 0.5+rnd32())

				hitR.delete = true

				// split rock
				nr := r.split(missile.getData().Pos, missile.getData().Speed, 6)

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
