package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type motion struct {
	pos, speed    V2
	rot, rotSpeed float64
}

func (m *motion) Move(dt float64) {
	dv := V2MulA(m.speed, dt)
	m.pos.Incr(dv)
	m.rot += m.rotSpeed * dt
}

type shape struct {
	points []V2
	bRect  Rect
}

func newShape(p []V2) *shape { return &shape{points: p} }

func (s *shape) Draw(pos V2, rot float64, colFill, colLine rl.Color) {
	var veryfirst, pp V2
	rotM := newM22rot(rot)

	var minx, maxx, miny, maxy int32
	for i, p := range s.points {
		np := rotM.pMulV(p)
		np.Incr(pos)
		if i > 0 {
			_triangle(np, pp, pos, colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn
			_line(pp, np, colLine)
			minx = min(minx, int32(np.x))
			maxx = max(maxx, int32(np.x))
			miny = min(miny, int32(np.y))
			maxy = max(maxy, int32(np.y))
		} else {
			veryfirst = np
			minx = int32(np.x)
			miny = int32(np.y)
			maxx, maxy = minx, miny
		}
		pp = np
	}
	_triangle(veryfirst, pp, pos, colFill)
	_line(pp, veryfirst, colLine)
	s.bRect = Rect{minx, miny, maxx - minx, maxy - miny}
}

func (s *shape) DrawThin(pos V2, rot float64, colFill, colLine rl.Color, thickness float32) {
	var veryfirst, pp V2
	rotM := newM22rot(rot)

	var minx, maxx, miny, maxy int32
	for i, p := range s.points {
		np := rotM.pMulV(p)
		np.Incr(pos)

		if i > 0 {
			_triangle(np, pp, pos, colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn
			_lineThick(pp, np, thickness, colLine)
			minx = min(minx, int32(np.x))
			maxx = max(maxx, int32(np.x))
			miny = min(miny, int32(np.y))
			maxy = max(maxy, int32(np.y))
		} else {
			veryfirst = np
			minx = int32(np.x)
			miny = int32(np.y)
			maxx, maxy = minx, miny
		}
		pp = np
	}

	_triangle(veryfirst, pp, pos, colFill)
	_lineThick(pp, veryfirst, thickness, colLine)
	s.bRect = Rect{minx, miny, maxx - minx, maxy - miny}
}
