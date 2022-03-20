package main

import "math"

type M22 struct {
	a00, a01, a10, a11 float64
}

func (m *M22) Mul(v V2) V2 {
	return M22MulV(*m, v)
}

type V2 struct {
	x, y float64
}

func (v1 *V2) Incr(v2 V2) {
	v1.x += v2.x
	v1.y += v2.y
}
func (v1 *V2) Decr(v2 V2) {
	v1.x -= v2.x
	v1.y -= v2.y
}

func V2MulA(v V2, a float64) V2 { // vector x scalar multiplication
	var nv V2
	nv.x, nv.y = v.x*a, v.y*a
	return nv
}
func (v1 V2) MulA(a float64) V2 {
	return V2{v1.x * a, v1.y * a}
}

func (v V2) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}
func (v V2) Len2() float64 { // Lenght suared
	return v.x*v.x + v.y*v.y
}
func V2Mul(v1, v2 V2) V2   { return V2{v1.x * v2.x, v1.y * v2.y} }
func (v1 V2) Mul(v2 V2) V2 { return V2{v1.x * v2.x, v1.y * v2.y} }

func V2Add(a, b V2) V2   { return V2{a.x + b.x, a.y + b.y} }
func (a V2) Add(b V2) V2 { return V2{a.x + b.x, a.y + b.y} }

func V2Sub(a, b V2) V2   { return V2{a.x - b.x, a.y - b.y} }
func (a V2) Sub(b V2) V2 { return V2{a.x - b.x, a.y - b.y} }

func (a V2) DivA(b float64) V2 { return V2{a.x / b, a.y / b} }

func (v V2) Norm() V2 {
	return v.DivA(v.Len())
}

func M22MulV(m M22, v V2) V2 { // matrix x vector multiplication
	var r V2
	r.x = m.a00*v.x + m.a01*v.y
	r.y = m.a10*v.x + m.a11*v.y
	return r
}
func M22pMulV(m *M22, v V2) V2 { // matrix x vector multiplication
	var r V2
	r.x = m.a00*v.x + m.a01*v.y
	r.y = m.a10*v.x + m.a11*v.y
	return r
}

func newM22Id() M22 {
	var m M22
	m.a00, m.a11 = 1, 1
	return m
}
func V2len(v V2) float64 { // Vector length
	return math.Sqrt(v.x*v.x + v.y*v.y)
}
func V2len2(v V2) float64 { // Vector length squared
	return v.x*v.x + v.y*v.y
}

func newM22rot(alpha float64) M22 {
	var m M22
	rad := float64(alpha * math.Pi / 180.0)
	m.a00, m.a01 = math.Cos(rad), -math.Sin(rad)
	m.a10, m.a11 = math.Sin(rad), math.Cos(rad)
	return m
}
func cs(alpha float64) V2 {
	rad := alpha * math.Pi / 180.0
	return V2{-math.Sin(rad), math.Cos(rad)}
}
