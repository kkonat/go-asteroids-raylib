package main

import (
	"math"
	"math/rand"
	qt "rlbb/lib/quadtree"
	v "rlbb/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rock struct {
	*shape
	motion
	radius float64
	mass   float64
	delete bool
}

func (r *Rock) BRect() qt.Rect {
	x := int32(r.pos.X)
	y := int32(r.pos.Y)
	rad := int32(r.radius)
	return qt.Rect{X: x - rad, Y: y - rad, W: rad * 2, H: rad * 2}
}

func (r *Rock) randomize() {
	r.radius = PrefferredRockSize/10 + rnd()*PrefferredRockSize
	r.mass = squared(r.radius)
	r.buildShape()
	r.rotSpeed = rnd()*1.5 - 0.75
}

func newRockRandom(g *game) *Rock {
	r := new(Rock)
	r.randomize()
	r.pos = V2{X: rnd() * g.gW, Y: rnd() * g.gH}
	r.speed = V2{X: rnd()*rSpeedMax*2.0 - rSpeedMax, Y: rnd()*rSpeedMax*2.0 - rSpeedMax}
	r.rotSpeed = rnd()*0.2 - 0.1
	return r
}

func newRockAt(pos, speed V2) *Rock {
	r := new(Rock)
	r.randomize()
	r.pos = pos
	r.speed = speed
	return r
}

func (r *Rock) buildShape() {
	n := 6 + rand.Intn(10) + int(r.radius/5)
	var step = 360 / float64(n)
	data := make([]V2, n)
	angle := 0.0
	for i := 0; i < n; i++ {
		angle += step + rnd()*step/2 - step/4
		r1 := r.radius + rnd()*r.radius/4 - r.radius/8
		p := v.RotV(angle)
		data[i] = p.MulA(r1)
	}
	r.shape = newShape(data)
}

func (g *game) generateRocks(preferredRocks int) {

	const safeCircle = 300
	cX, cY := g.ship.pos.X, g.ship.pos.Y
	i := 0
	for i < preferredRocks { // ( cx +r )  ( nr.X +nr.r)
		nr := newRockRandom(g)
		if cX+safeCircle < nr.pos.X+nr.radius || cX-safeCircle > nr.pos.X-nr.radius ||
			cY+safeCircle < nr.pos.Y+nr.radius || cY-safeCircle > nr.pos.Y-nr.radius {
			g.rocks.PushBack(nr)
			i++
		}
	}

}
func (r *Rock) Draw() { r.shape.DrawThin(r.pos, r.rot, rl.Black, rl.White, 0.75) }

func (r *Rock) split(hitat, speed V2, n int) []*Rock {

	rockDist2 := func(c1, c2 *Rock) float64 {
		return (c1.pos.X-c2.pos.X)*(c1.pos.X-c2.pos.X) +
			(c1.pos.Y-c2.pos.Y)*(c1.pos.Y-c2.pos.Y)
	}
	touches := func(which int, allRocks []*Rock) (bool, int) {
		for j, rock := range allRocks {
			if which != j {
				if rockDist2(allRocks[which], rock) < squared(allRocks[which].radius+rock.radius) {
					return true, j
				}
			}
		}
		return false, 0
	}

	// make new *Rocks and helper slices
	newRocks := make([]*Rock, n)
	frozen := make([]bool, n)

	// generate n random points: generate in polar coordinates convert to xy
	alphastep := math.Pi * 2.0 / float64(n)
	alpha := 0.0
	torim := r.radius
	for i := 0; i < n; i++ {
		dist := r.radius*0.75 - rnd()*r.radius/4
		torim = min(torim, r.radius-dist)
		newRocks[i] = newRockAt(V2{X: r.pos.X + math.Sin(alpha)*dist, Y: r.pos.Y + math.Cos(alpha)*dist}, V2{})
		frozen[i] = false
		alpha += alphastep
	}
	// find minimum distance between all points, initialy the minimum value is big circle radius,

	var mindist2 = squared(r.radius) // compare suared values to avoid square root

	for i, c1 := range newRocks {
		for j, c2 := range newRocks {
			if i != j {
				dist2 := rockDist2(c1, c2)
				if dist2 < mindist2 {
					mindist2 = dist2
				}
			}
		}
	}
	// circle initial radius = half that distance
	d := min(torim, math.Sqrt(mindist2)/2)

	// seed circles on all these points with this min radius - random value
	// they do not overlap
	for i := range newRocks {
		newRocks[i].radius = rnd() * d
	}

	// reposition new rocks
	for {
		var increased = 0
		for i := range newRocks {
			if !frozen[i] { // repeat until nothing moves
				d := math.Sqrt(rockDist2(newRocks[i], r))
				if d+newRocks[i].radius < r.radius { // until it touches the outer rock
					newRocks[i].radius += rnd() // grow radius
					increased++
				} else {
					//slide towards centre
					v := newRocks[i].pos.Sub(r.pos)
					newRocks[i].pos.Decr(v.DivA(d))

					d := math.Sqrt(rockDist2(newRocks[i], r))
					if d+newRocks[i].radius > r.radius { // if touches the outher rock
						frozen[i] = true
					}
				}

				t, j := touches(i, newRocks)
				if t {
					frozen[i] = true
					d = math.Sqrt(rockDist2(newRocks[i], newRocks[j]))
					if d > 0 {
						dd := newRocks[i].pos.Sub(newRocks[j].pos) // distance between i,j
						newRocks[i].pos.Incr(dd.DivA(d / 2))       // move away rock[i] from rock[j]
						newRocks[j].pos.Decr(dd.DivA(d / 2))

						d = math.Sqrt(rockDist2(newRocks[i], r))
						if d+newRocks[i].radius > r.radius {
							frozen[i] = true
						}
						d = math.Sqrt(rockDist2(newRocks[j], r))
						if d+newRocks[j].radius > r.radius {
							frozen[j] = true
						}
					}
				}
			}

		}
		if increased == 0 {
			break
		}
	}
	center := V2{} // calculate new center of the rocks group and calculate masses
	for i, r := range newRocks {
		r.radius *= 1.2
		newRocks[i].mass = squared(r.radius)
		center.Incr(r.pos)
	}
	center = center.DivA(float64(len(newRocks)))

	for i, ir := range newRocks {
		explodev := ir.pos.Sub(center).Norm()     // force throwing rocks outside
		rotv := V2{X: -explodev.Y, Y: explodev.X} //perpendicular
		rotv = rotv.MulA(r.rotSpeed * 5)          // centrifugal force
		explthrust := r.pos.Sub(hitat).Norm()     // thrust from the hit point
		missilethr := speed.Norm()                // missile speedd contribution
		masscontrib := math.Sqrt(ir.mass) / 5     // mass impact, larger move less

		newspeed := explodev.Add(rotv).Add(explthrust).Add(missilethr).DivA(masscontrib) // compute new rock speed
		newRocks[i].speed = r.speed.Add(newspeed)                                        // add it
		newRocks[i].rot = (r.rot + rnd()*2.0 - 1) / 2                                    // compute new rotation
	}
	return newRocks
}
