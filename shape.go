package main

import (
	qt "rlbb/lib/quadtree"
	v "rlbb/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type motion struct {
	pos, speed    V2
	rot, rotSpeed float64
}

func (m *motion) Move(dt float64) {
	dv := m.speed.MulA(dt)
	m.pos.Incr(dv)
	m.rot += m.rotSpeed * dt
}

type shape struct {
	points []V2
	bRect  qt.Rect
}

func newShape(p []V2) *shape {
	var s *shape
	s = &shape{points: p}
	s.points = append(s.points, s.points[0])
	return s

}

func (s *shape) Draw(pos V2, rot float64, colFill, colLine rl.Color) {
	var p1 V2
	rotM := v.NewM22rot(rot)

	var minx, maxx, miny, maxy int32
	for i, p := range s.points {
		p2 := rotM.PMulV(p)
		p2.Incr(pos)
		if i > 0 {
			_triangle(p2, p1, pos, colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn

			n1 := p1.Sub(pos).Norm()
			n2 := p2.Sub(pos).Norm()

			color := Game.Lights.ComputeColor(p1, p2, n1, n2, _ColorfromRlColor(colLine))

			_line(p1, p2, _rlColorFromColor(color))

			minx = min(minx, int32(p2.X))
			maxx = max(maxx, int32(p2.X))
			miny = min(miny, int32(p2.Y))
			maxy = max(maxy, int32(p2.Y))
		} else {
			//veryfirst = p2
			minx = int32(p2.X)
			miny = int32(p2.Y)
			maxx, maxy = minx, miny
		}
		p1 = p2
	}
	// 	_triangle(veryfirst, p1, pos, colFill)

	// 	_line(p1, veryfirst, colLine)
	s.bRect = qt.Rect{X: minx, Y: miny, W: maxx - minx, H: maxy - miny}
}

func (s *shape) DrawThin(pos V2, rot float64, colFill, colLine rl.Color, thickness float32) {
	var p1 V2
	rotM := v.NewM22rot(rot)

	var minx, maxx, miny, maxy int32

	for i, p := range s.points {
		p2 := rotM.PMulV(p)
		p2.Incr(pos)

		if i > 0 {
			_triangle(p2, p1, pos, colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn

			n1 := p1.Sub(pos).Norm()
			n2 := p2.Sub(pos).Norm()

			color := Game.Lights.ComputeColor(p1, p2, n1, n2, _ColorfromRlColor(colLine))

			// _line(p1, p1.Add(n1.MulA(5)), rl.White) // Debug draw nrormals
			// _line(p2, p2.Add(n2.MulA(5)), rl.White)

			_lineThick(p1, p2, thickness, _rlColorFromColor(color))

			minx = min(minx, int32(p2.X))
			maxx = max(maxx, int32(p2.X))
			miny = min(miny, int32(p2.Y))
			maxy = max(maxy, int32(p2.Y))
		} else {
			//veryfirst = p2
			minx = int32(p2.X)
			miny = int32(p2.Y)
			maxx, maxy = minx, miny
		}
		p1 = p2
	}

	// _triangle(veryfirst, p1, pos, colFill)
	// _lineThick(p1, veryfirst, thickness, colLine)
	s.bRect = qt.Rect{X: minx, Y: miny, W: maxx - minx, H: maxy - miny}
}
