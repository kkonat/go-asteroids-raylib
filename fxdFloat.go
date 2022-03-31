package main

type fxdFloat int32

const fxdFloatShift = 8

func FloatToFxdInt(a float64) fxdFloat {
	return fxdFloat(float64(1<<fxdFloatShift) * a)
}

type fxdV2 struct {
	x, y fxdFloat
}

func (v V2) ToV2int() fxdV2 {
	return fxdV2{FloatToFxdInt(v.x), FloatToFxdInt(v.y)}
}
func (v fxdV2) MulA(a int32) fxdV2 {
	return fxdV2{(v.x * fxdFloat(a)) >> fxdFloatShift, (v.y * fxdFloat(a)) >> fxdFloatShift}
}
func (v1 fxdV2) Add(v2 fxdV2) fxdV2 {
	return fxdV2{v1.x + v2.x, v1.y + v2.y}
}

type M22 struct {
	a00, a01, a10, a11 float64
}
