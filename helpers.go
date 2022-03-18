package main

import rl "github.com/gen2brain/raylib-go/raylib"

func _line(p1, p2 V2, col rl.Color) {
	rl.DrawLine(int32(p1.x), int32(p1.y), int32(p2.x), int32(p2.y), col)
}
func _circle(p1 V2, r float64, col rl.Color) {
	rl.DrawCircleLines(int32(p1.x), int32(p1.y), float32(r), col)
}
func _triangle(p1, p2, p3 V2, col rl.Color) {
	rl.DrawTriangle(rlV2(p1), rlV2(p2), rlV2(p3), col)
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
