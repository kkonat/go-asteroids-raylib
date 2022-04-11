package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Color struct {
	R, G, B, A float64
}

func (c *Color) Incr(col Color) {
	c.R += col.R
	c.G += col.G
	c.B += col.B
}
func (c1 Color) Add(c2 Color) Color { return Color{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, c1.A + c2.A} }
func (c1 Color) Mul(c2 Color) Color { return Color{c1.R * c2.R, c1.G * c2.G, c1.B * c2.B, c1.A * c2.A} }
func (c Color) DivA(d float64) Color {
	return Color{c.R / d, c.G / d, c.B / d, c.A / d}
}
func (c Color) MulA(m float64) Color {
	return Color{c.R * m, c.G * m, c.B * m, c.A * m}
}
func (c1 Color) Lerp(t float64, c2 Color) Color {
	return c1.MulA(1 - t).Add(c2.MulA(t))
}

func newColorRGB(r, g, b float64) Color { return Color{r, g, b, 1.0} }
func newColorRGBint(r, g, b uint8) Color {
	return Color{float64(r / 255), float64(g / 255), float64(b / 255), 1.0}
}
func nreColorBlack() Color { return Color{0, 0, 0, 1} }
func ColorToRlColor(c Color) rl.Color {
	return rl.Color{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), 255}
}
func ColorToRlColorFade(c Color, fade float64) rl.Color {
	return rl.Color{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), uint8(fade * 255)}
}

type LightSource interface {
	ComputeColor(at, normal V2, col Color) Color
	Draw()
	GetPos() V2
}

type Lighting struct {
	lights []LightSource
}

type OmniLight struct {
	Pos      V2
	Col      Color
	Sterngth float64
}
type SpotLight struct {
	OmniLight
	Dir   V2
	Angle float64
	Blur  float64
}

func newLighting() *Lighting {
	l := new(Lighting)
	return l
}

func (l *Lighting) AddLight(light LightSource) int {
	l.lights = append(l.lights, light)
	return len(l.lights) - 1
}
func (l *Lighting) DeleteLight(i int) {
	l.lights = append(l.lights[:i], l.lights[i+1:]...)
}

func (l Lighting) ComputeColor(p1, p2, n1, n2 V2, ownCol Color) Color {
	var col Color
	var i int

	n := n1.Add(n2).MulA(0.5)
	p := p1.Add(p2).MulA(0.5)
	for _, l := range l.lights {
		area := ((p2.X-p1.X)*(l.GetPos().Y-p1.Y) - (l.GetPos().X-p1.X)*(p2.Y-p1.Y)) // check if line segment is cw or ccw relative to light source
		if area < 0 {
			col.Incr(l.ComputeColor(p, n, ownCol))
			i = i + 1
		}
	}

	if i > 0 {
		return col.DivA(float64(i))
	} else {
		return newColorRGB(0, 0, 0) // invisible
	}

}
func (l *Lighting) Draw() {
	for _, l := range l.lights {
		l.Draw()
	}
}
func (ol *OmniLight) SetColor(col Color) { ol.Col = col }
func (ol OmniLight) GetPos() V2         { return ol.Pos }
func (ol OmniLight) ComputeColor(at, normal V2, col Color) Color {
	col = col.Mul(ol.Col)
	v := ol.Pos.Sub(at)
	dist := v.Len()
	l := v.DivA(dist)

	diffuse := normal.NormDot(l)
	diffuse = math.Abs(diffuse)
	diffuse *= 500 * ol.Sterngth / (dist*dist + 0.001)

	if diffuse > 1 {
		diffuse = 1
	}

	c := col.MulA(diffuse)
	return c
}
func (ol OmniLight) Draw() {
	rl.DrawCircleGradient(int32(ol.Pos.X), int32(ol.Pos.Y), float32(ol.Sterngth), ColorToRlColorFade(ol.Col, 0.2), rl.NewColor(0, 0, 0, 0))
}
func (ol *SpotLight) ComputeColor(at, normal V2, col rl.Color) rl.Color {
	return rl.Black
}
