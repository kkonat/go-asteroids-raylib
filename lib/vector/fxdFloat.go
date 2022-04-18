package vector

// simple fixed floating point math

type fxdFP int32 // basic data type

const fxdFloatShift = 8 // numver of fractional bits

// converts float64 to fixed floating point
func FloatToFxdFP(a float64) fxdFP {
	return fxdFP(float64(1<<fxdFloatShift) * a)
}

// fixed floating point 2D vector
type FxdFPV2 struct {
	X, Y fxdFP
}

// V2 to Fixed Floating point 2D vector conversion
func (vect V2) ToFxdFPV2() FxdFPV2 {
	return FxdFPV2{FloatToFxdFP(vect.X), FloatToFxdFP(vect.Y)}
}

// scalar mulatiplication
func (v FxdFPV2) MulA(a int32) FxdFPV2 {
	return FxdFPV2{(v.X * fxdFP(a)) >> fxdFloatShift, (v.Y * fxdFP(a)) >> fxdFloatShift}
}

// scalar division
func (v FxdFPV2) DivA(a int32) FxdFPV2 {
	return FxdFPV2{(v.X / fxdFP(a)) >> fxdFloatShift, (v.Y / fxdFP(a)) >> fxdFloatShift}
}

// vector addition
func (v1 FxdFPV2) Add(v2 FxdFPV2) FxdFPV2 {
	return FxdFPV2{v1.X + v2.X, v1.Y + v2.Y}
}

// converts fixedpoint floating value to int32
func (v fxdFP) ToInt32() int32 {
	return int32(v >> fxdFloatShift)
}
