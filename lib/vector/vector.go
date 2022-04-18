package vector

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type M22 struct {
	a00, a01, a10, a11 float64
}

func (m *M22) Mul(v V2) V2 {
	return m.MulV(v)
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

func (v1 V2) MulA(a float64) V2 {
	return V2{v1.X * a, v1.Y * a}
}

// Vector lenght
func (v V2) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Vector lenght suared
func (v V2) Len2() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (a V2) Add(b V2) V2       { return V2{a.X + b.X, a.Y + b.Y} }
func (a V2) Sub(b V2) V2       { return V2{a.X - b.X, a.Y - b.Y} }
func (a V2) SubA(b float64) V2 { return V2{a.X - b, a.Y - b} }
func (a V2) DivA(b float64) V2 { return V2{a.X / b, a.Y / b} }
func (a V2) PerpCCW() V2       { return V2{a.Y, -a.X} }
func (a V2) PerpCW() V2        { return V2{-a.Y, a.X} }
func RotV(angle float64) V2 {
	rad := angle * rl.Deg2rad
	return V2{math.Cos(rad), math.Sin(rad)}
}
func (v V2) Norm() V2 {
	return v.DivA(v.Len())
}
func (v1 V2) NormDot(v2 V2) float64 { return v1.X*v2.X + v1.Y*v2.Y }

// matrix by vector multiplication
func (m *M22) MulV(v V2) V2 {
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

func NewM22rot(alpha float64) M22 {
	// var m M22
	rad := float64((alpha - 90) * rl.Deg2rad)
	c, s := math.Cos(rad), math.Sin(rad)
	// m.a00, m.a01 = c, -s
	// m.a10, m.a11 = s, c
	return M22{c, -s, s, c}

}
