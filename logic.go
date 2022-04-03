package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *game) constrainShip() {
	const limit = 40.0
	const getback = 0.5
	if g.ship.pos.x < limit || g.ship.pos.x > g.gW-limit || g.ship.pos.y < limit || g.ship.pos.y > g.gH-limit {
		if g.ship.isSliding {
			g.ship.speed = V2MulA(g.ship.speed, 0.9)
			if g.ship.pos.Len2() < 0.01 {
				g.ship.isSliding = false
			}
		} else {
			if g.ship.pos.x < limit {
				g.ship.speed.x = getback
			}
			if g.ship.pos.x > g.gW-limit {
				g.ship.speed.x = -getback
			}
			if g.ship.pos.y < limit {
				g.ship.speed.y = getback
			}
			if g.ship.pos.y > g.gH-limit {
				g.ship.speed.y = -getback
			}
		}

	}
}

func (g *game) constrainRocks() {
	const limit = -50 // ordinate to detect rock going off screen / respawn new one

	iterator := g.rocks.Iter()
	for rock, ok := iterator(); ok; rock, ok = iterator() {
		r := rock.Value
		// for i := 0; i < len(g.rocks); i++ {	// TOTO: del

		if (r.pos.x+r.radius < limit && r.speed.x < 0) ||
			(r.pos.x-r.radius > g.gW-limit && r.speed.x > 0) ||
			(r.pos.y+r.radius < limit && r.speed.y < 0) ||
			(r.pos.y-r.radius > g.gH-limit && r.speed.y > 0) {
			if g.rocks.Len < noPreferredRocks {
				// respawn rock in a new sector, first randomly on the screen
				r.randomize()
				r.speed = V2{rnd()*rSpeedMax - rSpeedMax/2, rnd()*rSpeedMax - rSpeedMax/2}
				r.pos = V2{rnd() * g.gW, rnd() * g.gH}

				p := &r.pos     // these variables addes to make the
				s := &r.speed   // switch statement below more redable
				rad := r.radius // kind of hack'ish but reads easier

				sect := rand.Intn(4) // random sector from wchich new rock is to originate

				switch sect { // move the rock off the screen
				case 0: // left
					{
						p.x = -rad + limit //
						if s.x < 0 {
							s.x = -s.x
						}
					}
				case 1: // top
					{
						p.y = -rad + limit
						if s.y < 0 {
							s.x = -s.x
						}
					}
				case 2: // right
					{
						p.x = g.gW + rad - limit
						if s.x > 0 {
							s.x = -s.x
						}
					}
				case 3: // down
					{
						p.y = g.gH + rad - limit
						if s.y > 0 {
							s.y = -s.y
						}
					}
				}
			} else {
				if rock == nil || g.rocks.Len == 0 {
					panic("delete rock outside")
				}
				g.rocks.Delete(rock)
			}
		}
	}
}
func (g *game) constrainMissiles() {
	const limit = -10
	var i int
	for i < len(g.missiles) {
		p := g.missiles[i].pos
		if p.x <= limit || p.x > g.gW-limit ||
			p.y <= limit || p.y > g.gH-limit {

			g.deleteMissile(i)

		} else {
			i++
		}
	}
}

func (g *game) deleteMissile(m int) {
	if m > len(g.missiles) || m < 0 { // DEBUG
		panic("will crash here")
	} else {
		g.missiles = append(g.missiles[:m], g.missiles[m+1:]...)
	}

}

type circle struct {
	rect Rect
	p    V2
	r    float64
}

// func (c *circle) bRect() Rect {
// 	return c.rect
// }

// func newCircle(x, y, r int) *circle {
// 	o := new(circle)
// 	o.p = V2{0, 0}
// 	o.r = float64(r)
// 	o.rect.x, o.rect.y = x, y
// 	o.rect.w, o.rect.h = r/2, r/2
// 	return o
// }
func newCircleV2(p V2, r float64) *circle {
	return &circle{
		Rect{int32(p.x - r), int32(p.y - r),
			int32(r * 2), int32(r * 2)},
		p, r}
}

func (g *game) process_ship_hits() {
	v := cs(g.ship.rot)

	var circles [3]*circle

	circles[0] = newCircleV2(g.ship.pos.Sub(v.MulA(10)), 7)
	circles[1] = newCircleV2(g.ship.pos.Add(v.MulA(2)), 5)
	circles[2] = newCircleV2(g.ship.pos.Add(v.MulA(10)), 3)

	for _, c := range circles {
		_circle(c.p, c.r, rl.Yellow)
	}

	shipBRect := g.ship.shape.bRect
	potCols := g.RocksQt.MayCollide(shipBRect)

	for _, ro := range potCols {
		r := ro.Value
		for _, c := range circles {
			dist2 := r.pos.Sub(c.p).Len2()
			if dist2 < squared(r.radius+c.r) {
				if g.ship.shields > 0.7 {
					g.ship.shields -= 0.7
					g.sm.playFor(sScratch, 80)
				} else {
					if !debug {
						g.addParticle(newSparks(g.ship.pos, g.ship.speed, 300, 260, 5, rl.White, rl.Red))
						if !g.ship.destroyed {
							g.sm.play(sExplodeShip)
						}
						g.ship.destroyed = true
					}
					//game_over()
				}
			} else {
				g.sm.stop(sScratch)
			}

		}
	}
}

func (g *game) process_missile_hits() {
	const mr = 10 // missile radius

	for mi, missile := range g.missiles {
		mp := missile.pos
		//ms := missile.speed.MulA(13)
		missileBRect := missile.shape.bRect
		potCols := g.RocksQt.MayCollide(missileBRect)
		if debug {
			if degubDrawMissileLines {
				for _, c := range potCols {
					_line(mp, c.Value.pos, rl.Red)
				}
			}
		}

		for _, ro := range potCols {
			r := ro.Value
			rp := r.pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(r.radius+mr) { // hit
				// explosion vFX
				distBonus := g.ship.pos.Sub(r.pos).Len() / 200
				score := 1 + int(100/r.radius*distBonus/3)
				g.ship.cash += score
				str := fmt.Sprintf("+%d", score)
				g.addParticle(newTextPart(missile.pos, missile.speed, str, 16, 2, 0, false, rl.Yellow, rl.Red))
				g.addParticle(newExplosion(missile.pos, missile.speed, 30, 0.5))
				g.addParticle(newSparks(missile.pos, missile.speed, 100, 100, 2.0, rl.Orange, rl.Red))

				// sound
				g.sm.playPM(sExpl, 0.5+rnd32())

				// split rock
				nr := r.split(missile.pos, missile.speed, 6)

				// copy new rocks
				for i := 0; i < len(nr); i++ {
					if nr[i].radius > 8 {
						if g.rocks.Len < maxRocks {
							g.rocks.AppendVal(nr[i])
							nr[i].buildShape()
						}
					}
				}

				g.rocks.Delete(&ro.ListEl)

				if mi >= len(g.missiles) {
					//panic("bad missile delete")
				} else {
					g.deleteMissile(mi)
				}
				break
			}
		}
	}
}
