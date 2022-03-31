package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
	ns "github.com/ojrac/opensimplex-go"
	"golang.org/x/exp/constraints"
)

var noise [256]V2

// puts opensimplex noise into array for fx animations
func _initNoise() {
	const scaleDown = 10.0
	const blendRange = 32

	n := ns.NewNormalized(rand.Int63())
	for i := 0; i < 256; i++ {
		noise[i] = V2{n.Eval2(float64(i)/scaleDown, 0), n.Eval2(float64(i)/scaleDown, 1)}
	}

	// blend start and end values in the array
	for i, j := 256-blendRange, 0; j < blendRange; i, j = i+1, j+1 {
		t := float64(j / blendRange)
		noise[i] = V2{noise[i].x*t + noise[j].x*(1-t), noise[i].y*t + noise[j].y*(1-t)}
		j++
	}
}
func _noise1D(index uint8) float64 {
	return noise[index].x
}
func _noise2D(index uint8) V2 {
	return noise[index]
}
func _line(p1, p2 V2, col rl.Color) {
	rl.DrawLine(int32(p1.x), int32(p1.y), int32(p2.x), int32(p2.y), col)
}
func _lineThick(p1, p2 V2, thickness float32, col rl.Color) {
	rl.DrawLineEx(rl.Vector2{X: float32(p1.x), Y: float32(p1.y)}, rl.Vector2{X: float32(p2.x), Y: float32(p2.y)}, thickness, col)
}
func _circle(p1 V2, r float64, col rl.Color) {
	rl.DrawCircleLines(int32(p1.x), int32(p1.y), float32(r), col)
}
func _disc(p1 V2, r float64, col rl.Color) {
	rl.DrawCircle(int32(p1.x), int32(p1.y), float32(r), col)
}
func _gradientdisc(p1 V2, r float64, col1, col2 rl.Color) {
	rl.DrawCircleGradient(int32(p1.x), int32(p1.y), float32(r), col1, col2)
}
func _triangle(p1, p2, p3 V2, col rl.Color) {
	rl.DrawTriangle(rlV2(p1), rlV2(p2), rlV2(p3), col)
}
func _rect(p1 V2, d int32, col rl.Color) {
	rl.DrawRectangle(int32(p1.x), int32(p1.y), d, d, col)
}
func _rectFxdV2(p1 fxdV2, d int32, col rl.Color) {
	rl.DrawRectangle(int32(p1.x>>fxdFloatShift), int32(p1.y>>fxdFloatShift), d, d, col)
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
	return uint8(float32(a)*t + float32(b)*(1.0-t))
}
func _colorBlend(t0, maxt uint8, col1, col2 rl.Color) rl.Color {
	t := float32(t0) / float32(maxt)
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

// convert V2 to raylib's Vector2
func rlV2(p V2) rl.Vector2 {
	return rl.Vector2{X: float32(p.x), Y: float32(p.y)}
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
