package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"reflect"

	"rlbb/lib/vector"

	rl "github.com/gen2brain/raylib-go/raylib"
	ns "github.com/ojrac/opensimplex-go"
	"golang.org/x/exp/constraints"
)

type V2 = vector.V2

var noise [256]V2

// puts opensimplex noise into array for fx animations
func _initNoise() {
	const scaleDown = 10.0
	const blendRange = 32

	n := ns.NewNormalized(rand.Int63())
	for i := 0; i < 256; i++ {
		noise[i] = V2{X: n.Eval2(float64(i)/scaleDown, 0), Y: n.Eval2(float64(i)/scaleDown, 1)}
	}

	// blend start and end values in the array
	for i, j := 256-blendRange, 0; j < blendRange; i, j = i+1, j+1 {
		t := float64(j / blendRange)
		noise[i] = V2{X: noise[i].X*t + noise[j].X*(1-t), Y: noise[i].Y*t + noise[j].Y*(1-t)}
		j++
	}
}

func _noise1D(index uint8) float64 {
	return noise[index].X
}
func _noise2D(index uint8) V2 {
	return noise[index]
}
func _line(p1, p2 V2, col rl.Color) {
	rl.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), col)
}
func _lineThick(p1, p2 V2, thickness float32, col rl.Color) {
	rl.DrawLineEx(rl.Vector2{X: float32(p1.X), Y: float32(p1.Y)}, rl.Vector2{X: float32(p2.X), Y: float32(p2.Y)}, thickness, col)
}
func _circle(p1 V2, r float64, col rl.Color) {
	rl.DrawCircleLines(int32(p1.X), int32(p1.Y), float32(r), col)
}
func _circleGradient(p1 V2, r float64, col1, col2 rl.Color) {
	rl.DrawCircleGradient(int32(p1.X), int32(p1.Y), float32(r), col1, col2)
}
func _disc(p1 V2, r float64, col rl.Color) {
	rl.DrawCircle(int32(p1.X), int32(p1.Y), float32(r), col)
}
func _gradientdisc(p1 V2, r float64, col1, col2 rl.Color) {
	rl.DrawCircleGradient(int32(p1.X), int32(p1.Y), float32(r), col1, col2)
}
func _triangle(p1, p2, p3 V2, col rl.Color) {
	rl.DrawTriangle(rlV2(p1), rlV2(p2), rlV2(p3), col)
}
func _rect(p1 V2, d int32, col rl.Color) {
	rl.DrawRectangle(int32(p1.X), int32(p1.Y), d, d, col)
}
func _rectFxdV2(p1 vector.FxdV2, d int32, col rl.Color) {
	rl.DrawRectangle(p1.XInt32(), p1.YInt32(), d, d, col)
}

// writes text with multiple Colors
// x, y - location, size - font size, variadic args: value, color, value, color....
// draws a part of the text after each value, color pair
func _multicolorText(x, y int32, size int32, args ...interface{}) int32 {

	var width int32
	var msg string
	var col color.RGBA
	//var colType = reflect.TypeOf(col)
	var value reflect.Value
	width = 0
	for _, arg := range args {
		switch reflect.ValueOf(arg).Kind() {
		case reflect.TypeOf(col).Kind():
			col = reflect.ValueOf(arg).Interface().(color.RGBA)
			msg = fmt.Sprintf("%v ", value)
			rl.DrawText(msg, int32(x+width), y, size, col)
			width += rl.MeasureText(msg, size)
		case reflect.String: // print only these types
			value = reflect.ValueOf(arg)
		case reflect.Int:
			value = reflect.ValueOf(arg)
		case reflect.Float32:
			value = reflect.ValueOf(arg)
		case reflect.Float64:
			value = reflect.ValueOf(arg)
		default:
			panic("not allowed") // rudimentary error checking
		}
	}
	return int32(width)
}
func lerp(t float32, a, b uint8) uint8 {
	return uint8(float32(a)*(1.0-t) + float32(b)*t)
}
func _colorBlend(t0, maxt uint8, col1, col2 rl.Color) rl.Color {
	t := float32(t0) / float32(maxt)
	return rl.Color{
		lerp(t, col1.R, col2.R),
		lerp(t, col1.G, col2.G),
		lerp(t, col1.B, col2.B),
		lerp(t, col1.A, col2.A)}
}
func _colorBlendFloat(t float32, col1, col2 rl.Color) rl.Color {
	return rl.Color{
		lerp(t, col1.R, col2.R),
		lerp(t, col1.G, col2.G),
		lerp(t, col1.B, col2.B),
		lerp(t, col1.A, col2.A)}
}
func _colorBlendA(t0 float64, col1, col2 rl.Color) rl.Color {
	t := float32(t0)
	return rl.Color{
		lerp(t, col1.R, col2.R),
		lerp(t, col1.G, col2.G),
		lerp(t, col1.B, col2.B),
		lerp(t, col1.A, col2.A)}
}

func _rlColorFromFloats(r, g, b float64) rl.Color {
	return color.RGBA{R: uint8(255 * r), G: uint8(255 * g), B: uint8(255 * b), A: 255}
}
func _rlColorFromColor(c Color) rl.Color {
	return color.RGBA{R: uint8(255 * c.R), G: uint8(255 * c.G), B: uint8(255 * c.B), A: 255}
}
func _ColorfromRlColor(c rl.Color) Color {
	return newColorRGB(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

// convert V2 to raylib's Vector2
func rlV2(p V2) rl.Vector2 {
	return rl.Vector2{X: float32(p.X), Y: float32(p.Y)}
}
func itoVec2(x, y int32) rl.Vector2 {
	return rl.Vector2{X: float32(x), Y: float32(y)}
}
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}
func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func rnd() float64              { return rand.Float64() }
func rnd32() float32            { return rand.Float32() }
func rndSym(v float64) float64  { return rnd()*v*2.0 - v }
func squared(a float64) float64 { return a * a }
