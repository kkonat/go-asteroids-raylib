package main

import (
	"fmt"
	"image/color"
	"log"
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
	var count int
	for rock, ok := iterator(); ok; rock, ok = iterator() {
		count++
		r := rock.Value
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
	fmt.Println("Rock count:", count, " rocklist len=", g.rocks.Len)
}
func (g *game) constrainMissiles() {
	const limit = -10
	var i int
	for i < len(g.missiles) {
		p := g.missiles[i].getData().pos
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
		log.Print("bad missile delete ")
	} else {
		g.missiles = append(g.missiles[:m], g.missiles[m+1:]...)
	}

}

type circle struct {
	rect Rect
	p    V2
	r    float64
}

func newCircleV2(p V2, r float64) *circle {
	return &circle{
		Rect{int32(p.x - r), int32(p.y - r),
			int32(r * 2), int32(r * 2)},
		p, r}
}
func (g *game) drawForceField() {
	if g.ship.forceField {
		c := uint8(0)
		_circleGradient(g.ship.pos, forceFieldRadius,
			rl.Fade(color.RGBA{0, 200, 200, 255}, 0.25),
			rl.Fade(color.RGBA{0, 240, 200, 255}, 0.0))
		for a := float64(10) * rnd(); a < 360+float64(10)*rnd(); a += 1 + 10*rnd() {
			p := cs(a).MulA(float64(forceFieldRadius) * rnd())
			_lineThick(g.ship.pos, g.ship.pos.Add(p), rand.Float32()*20+10,
				rl.Fade(color.RGBA{0, 100, 100, 127}, 0.05))
			c++
		}
	}
}
func (g *game) processForceField() {

	if g.ship.forceField {
		dSpeed := V2{0, 0}

		iterator := g.rocks.Iter()
		for r, ok := iterator(); ok; r, ok = iterator() {
			rock := r.Value
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
	v := cs(g.ship.rot)

	var circles [3]*circle

	circles[0] = newCircleV2(g.ship.pos.Sub(v.MulA(10)), 8)
	circles[1] = newCircleV2(g.ship.pos.Add(v.MulA(2)), 6)
	circles[2] = newCircleV2(g.ship.pos.Add(v.MulA(10)), 4)

	// for _, c := range circles {
	// 	_circle(c.p, c.r, rl.Yellow)
	// }

	shipBRect := g.ship.shape.bRect
	g.RocksQtMutex.RLock()
	potCols := g.RocksQt.MayCollide(shipBRect)
	g.RocksQtMutex.RUnlock()
	for _, ro := range potCols {
		r := ro.Value
		for _, c := range circles {
			dist2 := r.pos.Sub(c.p).Len2()
			if dist2 < squared(r.radius+c.r) {
				// _disc(c.p, c.r, rl.Red)
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

func (g *game) processMissileHits() {
	const mr = 10 // missile radius
	var hit bool

	mi := 0
	for mi < len(g.missiles) {
		missile := g.missiles[mi]

		mp := missile.getData().pos
		//ms := missile.speed.MulA(13)
		missileBRect := missile.getData().shape.bRect

		g.RocksQtMutex.RLock()
		potCols := g.RocksQt.MayCollide(missileBRect)
		g.RocksQtMutex.RUnlock()

		if debug {
			if degubDrawMissileLines {
				for _, c := range potCols {
					_line(mp, c.Value.pos, rl.Red)
				}
			}
		}

		hit = false
		for i := range potCols {
			ro := potCols[i]
			r := ro.Value
			rp := r.pos
			dist2 := rp.Sub(mp).Len2()
			if dist2 < squared(r.radius+mr) { // hit
				// explosion vFX
				distBonus := g.ship.pos.Sub(r.pos).Len() / 200
				score := 1 + int(100/r.radius*distBonus/3)
				g.ship.cash += score
				str := fmt.Sprintf("+%d", score)
				g.addParticle(newTextPart(missile.getData().pos, missile.getData().speed, str, 16, 2, 0, false, rl.Yellow, rl.Red))
				g.addParticle(newExplosion(missile.getData().pos, missile.getData().speed, 30, 0.5))
				g.addParticle(newSparks(missile.getData().pos, missile.getData().speed, 100, 100, 2.0, rl.Orange, rl.Red))

				// sound
				g.sm.playPM(sExpl, 0.5+rnd32())

				// split rock
				nr := r.split(missile.getData().pos, missile.getData().speed, 6)

				// copy new rocks
				for i := 0; i < len(nr); i++ {
					if nr[i].radius > 8 {
						if g.rocks.Len < maxRocks {
							g.rocks.AppendVal(nr[i])
							nr[i].buildShape()
						}
					}
				}

				g.rocks.Delete(&potCols[i].ListEl)
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
