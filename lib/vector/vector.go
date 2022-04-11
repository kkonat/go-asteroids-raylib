package vector

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type M22 struct {
	a00, a01, a10, a11 float64
}

func (m *M22) Mul(v V2) V2 {
	return M22MulV(*m, v)
}

type V2 struct {
	X, Y float64
}

func (v1 *V2) Incr(v2 V2) {
	v1.X += v2.X
	v1.Y += v2.Y
}
func (v1 *V2) Decr(v2 V2) {
	v1.X -= v2.X
	v1.Y -= v2.Y
}

func V2MulA(v V2, a float64) V2 { // vector X scalar multiplication
	var nv V2
	nv.X, nv.Y = v.X*a, v.Y*a
	return nv
}
func (v1 V2) MulA(a float64) V2 {
	return V2{v1.X * a, v1.Y * a}
}

func (v V2) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func (v V2) Len2() float64 { // Lenght suared
	return v.X*v.X + v.Y*v.Y
}
func V2Mul(v1, v2 V2) V2   { return V2{v1.X * v2.X, v1.Y * v2.Y} }
func (v1 V2) Mul(v2 V2) V2 { return V2{v1.X * v2.X, v1.Y * v2.Y} }

func V2Add(a, b V2) V2   { return V2{a.X + b.X, a.Y + b.Y} }
func (a V2) Add(b V2) V2 { return V2{a.X + b.X, a.Y + b.Y} }

func V2Sub(a, b V2) V2         { return V2{a.X - b.X, a.Y - b.Y} }
func (a V2) Sub(b V2) V2       { return V2{a.X - b.X, a.Y - b.Y} }
func (a V2) SubA(b float64) V2 { return V2{a.X - b, a.Y - b} }

func (a V2) DivA(b float64) V2 { return V2{a.X / b, a.Y / b} }

func RotV(angle float64) V2 {
	rad := angle * rl.Deg2rad
	return V2{math.Sin(rad), math.Cos(rad)}
}
func (v V2) Norm() V2 {
	return v.DivA(v.Len())
}
func (v1 V2) NormDot(v2 V2) float64 { return v1.X*v2.X + v1.Y*v2.Y }
func M22MulV(m M22, v V2) V2 { // matrix x vector multiplication
	var r V2
	r.X = m.a00*v.X + m.a01*v.Y
	r.Y = m.a10*v.X + m.a11*v.Y
	return r
}
func M22pMulV(m *M22, v V2) V2 { // matrix x vector multiplication
	var r V2
	r.X = m.a00*v.X + m.a01*v.Y
	r.Y = m.a10*v.X + m.a11*v.Y
	return r
}
func (m *M22) PMulV(v V2) V2 {
	var r V2
	r.X = m.a00*v.X + m.a01*v.Y
	r.Y = m.a10*v.X + m.a11*v.Y
	return r
}

func NewM22Id() M22 {
	var m M22
	m.a00, m.a11 = 1, 1
	return m
}
func V2len(v V2) float64 { // Vector length
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func V2len2(v V2) float64 { // Vector length squared
	return v.X*v.X + v.Y*v.Y
}

func NewM22rot(alpha float64) M22 {
	var m M22
	rad := float64(alpha * rl.Deg2rad)
	m.a00, m.a01 = math.Cos(rad), -math.Sin(rad)
	m.a10, m.a11 = math.Sin(rad), math.Cos(rad)
	return m
}
func Cs(alpha float64) V2 {
	rad := (alpha + 90) * rl.Deg2rad
	return V2{math.Cos(rad), math.Sin(rad)}
	//return V2{-math.Sin(rad), math.Cos(rad)}
}
