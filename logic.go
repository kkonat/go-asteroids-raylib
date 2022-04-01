package main

import (
	"fmt"
	"math/rand"
	"sync"

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

	for i := 0; i < len(g.rocks); i++ {

		if (g.rocks[i].pos.x+g.rocks[i].radius < limit && g.rocks[i].speed.x < 0) ||
			(g.rocks[i].pos.x-g.rocks[i].radius > g.gW-limit && g.rocks[i].speed.x > 0) ||
			(g.rocks[i].pos.y+g.rocks[i].radius < limit && g.rocks[i].speed.y < 0) ||
			(g.rocks[i].pos.y-g.rocks[i].radius > g.gH-limit && g.rocks[i].speed.y > 0) {
			mutex.Lock()
			if len(g.rocks) < noPreferredRocks {
				// respawn rock in a new sector, first randomly on the screen
				g.rocks[i].randomize()
				g.rocks[i].speed = V2{rnd()*rSpeedMax - rSpeedMax/2, rnd()*rSpeedMax - rSpeedMax/2}
				g.rocks[i].pos = V2{rnd() * g.gW, rnd() * g.gH}

				p := &g.rocks[i].pos   // these variables addes to make the
				s := &g.rocks[i].speed // switch statement below more redable
				r := g.rocks[i].radius   // kind of hack'ish but reads easier

				sect := rand.Intn(4) // random sector from wchich new rock is to originate

				switch sect { // move the rock off the screen
				case 0: // left
					{
						p.x = -r + limit //
						if s.x < 0 {
							s.x = -s.x
						}
					}
				case 1: // top
					{
						p.y = -r + limit
						if s.y < 0 {
							s.x = -s.x
						}
					}
				case 2: // right
					{
						p.x = g.gW + r - limit
						if s.x > 0 {
							s.x = -s.x
						}
					}
				case 3: // down
					{
						p.y = g.gH + r - limit
						if s.y > 0 {
							s.y = -s.y
						}
					}
				}
			} else {
				g.deleteRock(i)
			}
			mutex.Unlock()
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
			mutex.Lock()
			g.deleteMissile(i)
			mutex.Unlock()
		} else {
			i++
		}
	}
}

func (g *game) deleteMissile(m int) {

	g.missiles = append(g.missiles[:m], g.missiles[m+1:]...)
}
func (g *game) deleteRock(r int) {
	g.rocks = append(g.rocks[:r], g.rocks[r+1:]...)
}

var mutex sync.Mutex

type circle struct {
	rect Rect
	p    V2
	r    float64
}

func (c *circle) bRect() Rect {
	return c.rect
}

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

	// for _, c := range circles {
	// 	_circle(c.p, c.r, rl.Yellow)
	// }

	for _, r := range g.rocks {
		for _, c := range circles {
			dist2 := r.pos.Sub(c.p).Len2()
			if dist2 < squared(r.radius+c.r) {
				if g.ship.shields > 0.7 {
					g.ship.shields -= 0.7
					g.sm.playFor(sScratch, 80)
				} else {
					g.addParticle(newSparks(g.ship.pos, g.ship.speed, 300, 260, 5, rl.White, rl.Red))
					if !g.ship.destroyed {
						g.sm.play(sExplodeShip)
					}
					g.ship.destroyed = true

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
	for i, r := range g.rocks {
		for m := range g.missiles {
			mp := g.missiles[m].pos
			rp := r.pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(r.radius+mr) { // hit
				// explosion vFX
				distBonus := g.ship.pos.Sub(r.pos).Len() / 200
				score := 1 + int(100/r.radius*distBonus/3)
				g.ship.cash += score
				str := fmt.Sprintf("+%d", score)
				g.addParticle(newTextPart(g.missiles[m].pos, g.missiles[m].speed, str, 16, 2, 0, false, rl.Yellow, rl.Red))
				g.addParticle(newExplosion(g.missiles[m].pos, g.missiles[m].speed, 30, 0.5))
				g.addParticle(newSparks(g.missiles[m].pos, g.missiles[m].speed, 100, 100, 2.0, rl.Orange, rl.Red))

				// sound
				g.sm.playPM(sExpl, 0.5+rnd32())

				// split rock
				nr := g.rocks[i].split(g.missiles[m].pos, g.missiles[m].speed, 6)

				// copy new rocks
				for i := 0; i < len(nr); i++ {
					if nr[i].radius > 8 {
						if len(g.rocks) < maxRocks {
							g.rocks = append(g.rocks, nr[i])
							nr[i].buildShape()
						}
					}
				}
				g.deleteRock(i)
				g.deleteMissile(m)
				break
			}
		}
	}
}
