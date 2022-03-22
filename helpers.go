package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func _line(p1, p2 V2, col rl.Color) {
	rl.DrawLine(int32(p1.x), int32(p1.y), int32(p2.x), int32(p2.y), col)
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
func _square(p1 V2, d int32, col rl.Color) {
	rl.DrawRectangle(int32(p1.x), int32(p1.y), d, d, col)
}
func lerp(t float32, a, b uint8) uint8 {
	return uint8(float32(a)*t + float32(b)*(1.0-t))
}
func _colorBlend(t0, maxt int, col1, col2 rl.Color) rl.Color {
	t := float32(t0) / float32(maxt)
	return rl.Color{lerp(t, col1.R, col2.R),
		lerp(t, col1.G, col2.G),
		lerp(t, col1.B, col2.B),
		lerp(t, col1.A, col2.A)}

}
func rlV2(p V2) rl.Vector2 {
	return rl.Vector2{X: float32(p.x), Y: float32(p.y)}
}
func min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func rnd() float64              { return rand.Float64() }
func rnd32() float32            { return rand.Float32() }
func rndSym(v float64) float64  { return rnd()*v*2.0 - v }
func squared(a float64) float64 { return a * a }