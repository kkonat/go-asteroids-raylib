package main

import "math"

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
func V2len(v V2) float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y + v.y)))

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
