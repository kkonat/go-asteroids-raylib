package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type shape struct {
	points []V2

	pos, speed V2

	rot, rotSpeed float64
	rotM          M22
}

func newShape(p []V2) *shape {
	s := new(shape)
	s.points = p
	s.rotM = M22Id()
	return s
}
// func newEmptyShape() *shape {
// 	s := new(shape)
// 	s.rotM = M22Id()
// 	return s
// }

// func (s *shape) addPoints() {
// 	s.points = p
// }

func (s *shape) Draw(colFill, colLine rl.Color) {
	var ppx, ppy, npx, npy int32
	var v V2

	for i := 0; i < len(s.points); i++ {
		np := M22MulV(s.rotM, s.points[i])
		np = V2Add(np, s.pos)
		npx, npy = int32(np.x), int32(np.y)
		if i > 0 {
			rl.DrawTriangle( // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn
				rl.Vector2{X: float32(npx), Y: float32(npy)},
				rl.Vector2{X: float32(ppx), Y: float32(ppy)},
				rl.Vector2{X: float32(s.pos.x), Y: float32(s.pos.y)}, colFill)
			rl.DrawLine(ppx, ppy, npx, npy, colLine)
		} else {
			v = np
		}
		ppx, ppy = npx, npy
	}
	// v = M22V2mul(s.rotM, s.points[0])
	// np := V2add(s.pos, v)
	npx, npy = int32(v.x), int32(v.y)
	rl.DrawTriangle( // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn
		rl.Vector2{X: float32(npx), Y: float32(npy)},
		rl.Vector2{X: float32(ppx), Y: float32(ppy)},
		rl.Vector2{X: float32(s.pos.x), Y: float32(s.pos.y)}, colFill)
	rl.DrawLine(ppx, ppy, npx, npy, colLine)
	s.pos.Incr(s.speed)
	s.rotM = M22rot(s.rot)
	s.rot += s.rotSpeed
}

type trail struct {
	trail                []V2
	trailAge             []int
	trailHead, trailTail int
	color                rl.Color
	width                float32
}

func newTrail(size uint16, c rl.Color, w float32) *trail {
	t := new(trail)
	t.trail = make([]V2, size)
	t.trailAge = make([]int, size)
	t.color = c
	t.width = w
	return t
}
func (t *trail) addPoint(p V2) {
	idx := t.trailTail
	if idx < len(t.trail) {
		t.trail[idx] = p
		t.trailAge[idx]++
		t.trailTail++
	}
}
