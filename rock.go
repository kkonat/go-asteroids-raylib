package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rock struct {
	shape  *shape
	m      *motion
	radius float64
	mass   float64
}

func newRockRandom(g *game) *Rock {

	r := new(Rock)
	r.m = newMotion()
	r.randomize()
	r.m.pos = V2{rand.Float64() * float64(g.sW), rand.Float64() * float64(g.sH)}
	r.m.speed = V2{rand.Float64()*rSpeedMax*2.0 - rSpeedMax, rand.Float64()*rSpeedMax*2.0 - rSpeedMax}
	r.m.rotSpeed = rand.Float64()*0.2 - 0.1

	return r
}
func newRockAt(pos, speed V2) *Rock {

	r := new(Rock)
	r.m = newMotion()

	r.randomize()
	r.m.pos = pos
	r.m.speed = speed

	return r
}

func (r *Rock) buildShape() {
	n := 6 + rand.Intn(10) + int(r.radius/5)
	var step = 360 / float64(n)

	data := make([]V2, n)
	angle := 0.0
	for i := 0; i < n; i++ {
		angle += step + rand.Float64()*step/2 - step/4
		r1 := r.radius + rand.Float64()*r.radius/4 - r.radius/8
		p := cs(angle)
		data[i] = p.MulA(r1)
	}

	r.shape = newShape(data)

}

func (r *Rock) randomize() {
	r.radius = 10 + rand.Float64()*100
	//n := 6 + rand.Intn(10) + int(r.radius/5)

	r.buildShape()

	r.m.rotSpeed = rand.Float64()*1.5 - 0.75

}

func (r *Rock) Draw() {
	//rl.DrawCircle(int32(r.shape.pos.x), int32(r.shape.pos.y), r.radius, rl.ColorAlpha(rl.DarkGray, 0.2))
	//	rl.DrawLine(720,360,int32(r.shape.pos.x),int32(r.shape.pos.y),rl.DarkBlue)
	r.shape.Draw(r.m, rl.Black, rl.DarkGray)
}

func touches(which int, allRocks []*Rock) (bool, int) {
	for j, rock := range allRocks {
		if which != j {
			if dist2(allRocks[which], rock) < squared(allRocks[which].radius+rock.radius) {
				return true, j
			}
		}
	}
	return false, 0
}

func squared(a float64) float64 { return a * a }
func dist2(c1, c2 *Rock) float64 {
	return (c1.m.pos.x-c2.m.pos.x)*(c1.m.pos.x-c2.m.pos.x) +
		(c1.m.pos.y-c2.m.pos.y)*(c1.m.pos.y-c2.m.pos.y)
}
func (r *Rock) split(hitat, speed V2, n int) []*Rock {

	newRocks := make([]*Rock, n)
	frozen := make([]bool, n)

	// generate n random points: generate in polar coordinates convert to xy
	alphastep := float64(math.Pi*2) / float64(n)
	alpha := float64(0)
	torim := r.radius
	for i := 0; i < n; i++ {
		dist := r.radius*0.75 - rand.Float64()*r.radius/4
		torim = min(torim, r.radius-dist)
		newRocks[i] = newRockAt(V2{r.m.pos.x + math.Sin(alpha)*dist, r.m.pos.y + math.Cos(alpha)*dist}, V2{0, 0})
		frozen[i] = false
		alpha += alphastep
	}
	// find minimum distance between all points, initialy minimum is big circle radius,
	var mindist2 = squared(r.radius) // compare suared values to avoid square root

	for i, c1 := range newRocks {
		for j, c2 := range newRocks {
			if i != j {
				dist2 := dist2(c1, c2)
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
		newRocks[i].radius = rand.Float64() * d
	}
	//draw_circles(circles) // draw
	//rl.EndDrawing()

	for {

		var increased = 0
		for i := range newRocks {
			if !frozen[i] { // repeat until nothing moves
				d := math.Sqrt(dist2(newRocks[i], r))
				if d+newRocks[i].radius < r.radius {
					newRocks[i].radius += rand.Float64()
					increased++
				} else {
					//slide towards centre
					newRocks[i].m.pos.x -= (newRocks[i].m.pos.x - r.m.pos.x) / d
					newRocks[i].m.pos.y -= (newRocks[i].m.pos.y - r.m.pos.y) / d
					d := math.Sqrt(dist2(newRocks[i], r))
					if d+newRocks[i].radius > r.radius {
						frozen[i] = true
					}
				}

				t, j := touches(i, newRocks)
				if t {
					frozen[i] = true
					d = math.Sqrt(dist2(newRocks[i], newRocks[j]))
					if d > 0 {
						dx := (newRocks[i].m.pos.x - newRocks[j].m.pos.x) / d / 2 // vector along which it touches
						dy := (newRocks[i].m.pos.y - newRocks[j].m.pos.y) / d / 2
						newRocks[i].m.pos.x += dx
						newRocks[j].m.pos.x -= dx
						newRocks[i].m.pos.y += dy
						newRocks[j].m.pos.y -= dy
						d = math.Sqrt(dist2(newRocks[i], r))
						if d+newRocks[i].radius > r.radius {
							frozen[i] = true
						}
						d = math.Sqrt(dist2(newRocks[j], r))
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
	return newRocks
}
