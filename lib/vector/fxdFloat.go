package vector

type fxdFloat int32

const fxdFloatShift = 8

func FloatToFxdInt(a float64) fxdFloat {
	return fxdFloat(float64(1<<fxdFloatShift) * a)
}

type FxdV2 struct {
	X, Y fxdFloat
}

func (vect V2) ToV2int() FxdV2 {
	return FxdV2{FloatToFxdInt(vect.X), FloatToFxdInt(vect.Y)}
}
func (v FxdV2) MulA(a int32) FxdV2 {
	return FxdV2{(v.X * fxdFloat(a)) >> fxdFloatShift, (v.Y * fxdFloat(a)) >> fxdFloatShift}
}
func (v1 FxdV2) Add(v2 FxdV2) FxdV2 {
	return FxdV2{v1.X + v2.X, v1.Y + v2.Y}
}
func (p1 FxdV2) XInt32() int32 {
	return int32(p1.X >> fxdFloatShift)
}
func (p1 FxdV2) YInt32() int32 {
	return int32(p1.Y >> fxdFloatShift)
}
