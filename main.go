package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type M22 struct {
	a00, a01, a10, a11 float32
}

type V2 struct {
	x, y float32
}

func V2add(a, b V2) V2 {
	var v V2
	v.x, v.y = a.x+b.x, a.y+b.y
	return v
}
func V2V(a V2, b float32) V2 {
	var v V2
	v.x, v.y = a.x*b, a.y*b
	return v
}
func M22V2mul(m M22, v V2) V2 {
	var r V2
	r.x = m.a00*v.x + m.a01*v.y
	r.y = m.a10*v.x + m.a11*v.y
	return r
}
func M22Ident() M22 {
	var m M22
	m.a00, m.a11 = 1, 1
	return m
}

func M22rot(alpha float32) M22 {
	var m M22
	rad := float64(alpha * math.Pi / 180.0)
	m.a00, m.a01 = float32(math.Cos(rad)), -float32(math.Sin(rad))
	m.a10, m.a11 = float32(math.Sin(rad)), float32(math.Cos(rad))
	return m
}
func cs(alpha float32) V2 {
	rad := float64(alpha * math.Pi / 180.0)
	return V2{float32(-math.Sin(rad)), float32(math.Cos(rad))}
}

type vehicle interface {
	Draw()
}

type shape struct {
	points []V2
}

func (s *shape) add(p V2) {
	s.points = append(s.points, p)
}

type ship struct {
	pos, speed, ds V2

	rot  float32
	shp  shape
	rotM M22
	mass float32
	fuel float32
	col  rl.Color
}

func newShip(mass, fuel float32) *ship {
	const S = 16
	s := new(ship)
	s.shp.add(V2{-S / 2, -S})
	s.shp.add(V2{0, S})
	s.shp.add(V2{S / 2, -S})

	s.col = rl.White
	s.rotM = M22Ident()
	s.mass = mass
	s.fuel = fuel
	return s
}
func (s *ship) Draw() {
	var ppx, ppy, npx, npy int32
	var v V2
	for i, p := range s.shp.points {
		np := M22V2mul(s.rotM, p)
		np = V2add(np, s.pos)
		npx, npy = int32(np.x), int32(np.y)
		if i > 0 {
			rl.DrawLine(ppx, ppy, npx, npy, s.col)
		}
		ppx, ppy = npx, npy
	}
	v = M22V2mul(s.rotM, s.shp.points[0])
	np := V2add(s.pos, v)
	npx, npy = int32(np.x), int32(np.y)
	rl.DrawLine(ppx, ppy, npx, npy, s.col)
	s.pos = V2add(s.pos, s.speed)
	s.speed = V2V(s.speed, 0.997)

	rl.DrawLineEx(rl.Vector2{X: s.pos.x, Y: s.pos.y}, rl.Vector2{X: s.pos.x - s.ds.x*200, Y: s.pos.y - s.ds.y*200}, 4.1, rl.Orange)
}
func (s *ship) thrust(fuelCons float32) {
	//if s.fuel > 0 {
	force := fuelCons
	//	s.fuel -= fuelCons
	a := force * 200 / (s.mass + s.fuel)
	s.ds = V2V(cs(s.rot), a)
	s.speed = V2add(s.speed, s.ds)
}

func (s *ship) rotate(deltaAngle float32) {
	s.rot += deltaAngle
	s.rotM = M22rot(s.rot)
}

func main() {

	game := newGame(1440, 720)

	for !rl.WindowShouldClose() {
		if !game.sm.isPlaying(0) {
			game.sm.play(game.sm.sSpace)

			fmt.Println("started")
		}

		if rl.IsKeyPressed('Q') {
			game.sm.playM(game.sm.sOinx)
		}
		if rl.IsKeyPressed('M') {
			if !game.sm.mute {
				game.sm.stop(0)
			}
			game.sm.mute = !game.sm.mute

		}
		if rl.IsKeyDown('A') {
			game.ship.rotate(-2)
		}
		if rl.IsKeyDown('S') {
			game.ship.thrust(1.0)
			if !game.sm.isPlaying(game.sm.sThrust) {
				game.sm.play(game.sm.sThrust)
			}
		}
		if rl.IsKeyReleased('S') {
			game.ship.thrust(0)
			game.sm.stop(game.sm.sThrust)
		}
		if rl.IsKeyDown('D') {
			game.ship.rotate(2)
		}

		game.drawGame()

	}
	game.finalize()

}
