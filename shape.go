package main

import (
	qt "bangbang/lib/quadtree"
	v "bangbang/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type shape struct {
	points           []V2
	bRect            qt.Rect
	colFill, colLine rl.Color
}

func NewShape(p []V2, Fill, Line rl.Color) *shape {
	var s *shape
	s = &shape{points: p}
	s.points = append(s.points, s.points[0])
	s.colFill = Fill
	s.colLine = Line
	return s

}

func (s *shape) Draw(Pos v.V2, Rot float64) {
	var p1 V2
	rotM := v.NewM22rot(Rot)

	var minx, maxx, miny, maxy int32
	for i, p := range s.points {
		p2 := rotM.MulV(p)
		p2.Incr(Pos)
		if i > 0 {
			_triangle(p2, p1, Pos, s.colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn

			n1 := p1.Sub(Pos).Norm()
			n2 := p2.Sub(Pos).Norm()

			color := Game.VisibleLights.ComputeColor(p1, p2, n1, n2, _ColorfromRlColor(s.colLine))

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

func (s *shape) DrawThin(Pos v.V2, Rot float64, thickness float32) {
	var p1 V2
	rotM := v.NewM22rot(Rot)

	var minx, maxx, miny, maxy int32

	for i, p := range s.points {
		p2 := rotM.MulV(p)
		p2.Incr(Pos)

		if i > 0 {
			_triangle(p2, p1, Pos, s.colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn

			n1 := p1.Sub(Pos).Norm()
			n2 := p2.Sub(Pos).Norm()

			color := Game.VisibleLights.ComputeColor(p1, p2, n1, n2, _ColorfromRlColor(s.colLine))

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
