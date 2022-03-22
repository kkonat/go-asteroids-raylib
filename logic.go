package main

import (
	"math/rand"
	"sync"
)

func (g *game) constrainShip() {
	const limit = 40.0
	const getback = 0.5
	if g.ship.m.pos.x < limit || g.ship.m.pos.x > g.gW-limit || g.ship.m.pos.y < limit || g.ship.m.pos.y > g.gH-limit {
		if g.ship.isSliding {
			g.ship.m.speed = V2MulA(g.ship.m.speed, 0.9)
			if g.ship.m.pos.Len2() < 0.01 {
				g.ship.isSliding = false
			}
		} else {
			if g.ship.m.pos.x < limit {
				g.ship.m.speed.x = getback
			}
			if g.ship.m.pos.x > g.gW-limit {
				g.ship.m.speed.x = -getback
			}
			if g.ship.m.pos.y < limit {
				g.ship.m.speed.y = getback
			}
			if g.ship.m.pos.y > g.gH-limit {
				g.ship.m.speed.y = -getback
			}
		}

	}
}

func (g *game) constrainRocks() {
	const limit = -50 // ordinate to detect rock going off screen / respawn new one

	for i := 0; i < g.rocksNo; i++ {

		if (g.rocks[i].m.pos.x+g.rocks[i].radius < limit && g.rocks[i].m.speed.x < 0) ||
			(g.rocks[i].m.pos.x-g.rocks[i].radius > g.gW-limit && g.rocks[i].m.speed.x > 0) ||
			(g.rocks[i].m.pos.y+g.rocks[i].radius < limit && g.rocks[i].m.speed.y < 0) ||
			(g.rocks[i].m.pos.y-g.rocks[i].radius > g.gH-limit && g.rocks[i].m.speed.y > 0) {
			mutex.Lock()
			if g.rocksNo < preferredRocks {
				// respawn rock in a new sector, first randomly on the screen
				g.rocks[i].randomize()
				g.rocks[i].m.speed = V2{rnd()*rSpeedMax - rSpeedMax/2, rnd()*rSpeedMax - rSpeedMax/2}
				g.rocks[i].m.pos = V2{rnd() * g.gW, rnd() * g.gH}

				p := &g.rocks[i].m.pos   // these variables addes to make the
				s := &g.rocks[i].m.speed // switch statement below more redable
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

	i := 0
	for i < g.missilesNo {
		p := g.missiles[i].m.pos
		if p.x <= limit || p.x > g.gW-limit ||
			p.y <= limit || p.y > g.gH-limit {
			mutex.Lock()
			g.deleteMissile(i)
			mutex.Unlock()
			continue
		}
		i++
	}
}

func (g *game) deleteMissile(m int) {
	g.missilesNo--
	g.missiles[m] = g.missiles[g.missilesNo]
	g.missiles[g.missilesNo] = nil
}
func (g *game) deleteRock(r int) {
	g.rocksNo--
	g.rocks[r] = g.rocks[g.rocksNo]
	g.rocks[g.rocksNo] = nil
}

var mutex sync.Mutex

func (g *game) process_missile_hits() {
	const mr = 10 // missile radius
	for r := 0; r < g.rocksNo; r++ {
		for m := 0; m < g.missilesNo; m++ {
			mp := g.missiles[m].m.pos
			rp := g.rocks[r].m.pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(g.rocks[r].radius+mr) { // hit
				// explosion vFX
				g.addParticle(newExplosion(g.missiles[m].m.pos, g.missiles[m].m.speed, 30, 0.5))
				g.addParticle(newSparks(g.missiles[m].m.pos, g.missiles[m].m.speed, 100, 2.0))

				// sound
				g.sm.playPM(g.sm.sExpl, 0.5+rnd32())

				// split rock
				nr := g.rocks[r].split(g.missiles[m].m.pos, g.missiles[m].m.speed, 6)

				// copy new rocks
				for i := 0; i < len(nr); i++ {
					if nr[i].radius > 8 {
						if g.rocksNo < maxRocks {
							g.rocks[g.rocksNo] = nr[i]
							g.rocks[g.rocksNo].buildShape()
							g.rocksNo++
						}
					}
				}
				g.deleteRock(r)
				g.deleteMissile(m)
				break
			}
		}
	}
}
