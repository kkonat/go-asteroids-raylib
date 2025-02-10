package main

import (
	qt "bangbang/lib/quadtree"
	v "bangbang/lib/vector"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rock struct {
	*shape
	Motion
	radius float64
	mass   float64
	delete bool
}

func (r *Rock) BRect() qt.Rect {
	x := int32(r.Pos.X)
	y := int32(r.Pos.Y)
	rad := int32(r.radius)
	return qt.Rect{X: x - rad, Y: y - rad, W: rad * 2, H: rad * 2}
}

func (r *Rock) randomize() {
	r.radius = PrefferredRockSize/10 + rnd()*PrefferredRockSize
	r.mass = squared(r.radius)
	r.buildShape()
	r.RotSpeed = rnd()*1.5 - 0.75
}

func newRockRandom(g *game) *Rock {
	r := new(Rock)
	r.randomize()
	r.Pos = V2{X: rnd() * g.gW, Y: rnd() * g.gH}
	r.Speed = V2{X: rnd()*rSpeedMax*2.0 - rSpeedMax, Y: rnd()*rSpeedMax*2.0 - rSpeedMax}
	r.RotSpeed = rnd()*0.2 - 0.1
	return r
}

func newRockAt(pos, speed V2) *Rock {
	r := new(Rock)
	r.randomize()
	r.Pos = pos
	r.Speed = speed
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
	r.shape = NewShape(data, rl.Black, rl.White)
}

func (g *game) generateRocks(preferredRocks int) {

	const safeCircle = 300
	cX, cY := g.ship.Pos.X, g.ship.Pos.Y
	i := 0
	for i < preferredRocks { // ( cx +r )  ( nr.X +nr.r)
		nr := newRockRandom(g)
		if cX+safeCircle < nr.Pos.X+nr.radius || cX-safeCircle > nr.Pos.X-nr.radius ||
			cY+safeCircle < nr.Pos.Y+nr.radius || cY-safeCircle > nr.Pos.Y-nr.radius {
			g.rocks.PushBack(nr)
			i++
		}
	}

}
func (r *Rock) Draw() { r.shape.DrawThin(r.Pos, r.Rot, 0.75) }

func (r *Rock) split(hitat, speed V2, n int) []*Rock {

	rockDist2 := func(c1, c2 *Rock) float64 {
		return (c1.Pos.X-c2.Pos.X)*(c1.Pos.X-c2.Pos.X) +
			(c1.Pos.Y-c2.Pos.Y)*(c1.Pos.Y-c2.Pos.Y)
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
		newRocks[i] = newRockAt(V2{X: r.Pos.X + math.Sin(alpha)*dist, Y: r.Pos.Y + math.Cos(alpha)*dist}, V2{})
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
					v := newRocks[i].Pos.Sub(r.Pos)
					newRocks[i].Pos.Decr(v.DivA(d))

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
						dd := newRocks[i].Pos.Sub(newRocks[j].Pos) // distance between i,j
						newRocks[i].Pos.Incr(dd.DivA(d / 2))       // move away rock[i] from rock[j]
						newRocks[j].Pos.Decr(dd.DivA(d / 2))

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
		center.Incr(r.Pos)
	}
	center = center.DivA(float64(len(newRocks)))

	for i, ir := range newRocks {
		explodev := ir.Pos.Sub(center).Norm()     // force throwing rocks outside
		rotv := V2{X: -explodev.Y, Y: explodev.X} //perpendicular
		rotv = rotv.MulA(r.RotSpeed * 5)          // centrifugal force
		explthrust := r.Pos.Sub(hitat).Norm()     // thrust from the hit point
		missilethr := speed.Norm()                // missile speedd contribution
		masscontrib := math.Sqrt(ir.mass) / 5     // mass impact, larger move less

		newspeed := explodev.Add(rotv).Add(explthrust).Add(missilethr).DivA(masscontrib) // compute new rock speed
		newRocks[i].Speed = r.Speed.Add(newspeed)                                        // add it
		newRocks[i].Rot = (r.Rot + rnd()*2.0 - 1) / 2                                    // compute new rotation
	}
	return newRocks
}
