package main

import (
	"math"
	v "rlbb/lib/vector"

	"github.com/fogleman/ease"
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

// QuadOut easing
// t: current time, b: begInnIng value, c: change In value, d: duration
func easeOut(t float64) float64 {
	return ease.OutQuart(t)
}
func (c1 Color) gamma(gamma float64) Color {
	c1.R = math.Pow(c1.R, 1/gamma)
	c1.G = math.Pow(c1.G, 1/gamma)
	c1.B = math.Pow(c1.B, 1/gamma)
	return c1
}
func (c1 Color) clamp() Color {
	if c1.R > 1.0 {
		c1.R = easeOut(c1.R)
	}
	if c1.G > 1.0 {
		c1.G = easeOut(c1.R)
	}
	if c1.B > 1.0 {
		c1.B = easeOut(c1.R)
	}
	return c1
}

func newColorRGB(r, g, b float64) Color { return Color{r, g, b, 1.0} }
func newColorRGBint(r, g, b uint8) Color {
	return Color{float64(r) / 255, float64(g) / 255, float64(b) / 255, 1.0}
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
	sources []LightSource
}

var slshader rl.Shader

type OmniLight struct {
	Pos      V2
	Col      Color
	Strength float64
}
type SpotLight struct {
	OmniLight
	Dir   float64
	Angle float64
	Blur  float64
}

func newLighting() *Lighting {
	l := new(Lighting)
	l.sources = make([]LightSource, 0, 120)
	slshader = rl.LoadShader("shaders/base.vs", "shaders/spotlight.fs")
	return l
}

func (l *Lighting) AddLight(light LightSource) (lightIdx int) {
	l.sources = append(l.sources, light)
	lightIdx = len(l.sources) - 1
	return
}

func (l *Lighting) DeleteLight(light LightSource) {
	for li := range l.sources {
		if l.sources[li] == light {
			l.sources = append(l.sources[:li], l.sources[li+1:]...)
			return
		}
	}
	panic("can't delete PANIC: ")
}

func (l Lighting) ComputeColor(p1, p2, n1, n2 V2, ownCol Color) Color {
	var col Color
	//	var i int

	n := n1.Add(n2).MulA(0.5)
	p := p1.Add(p2).MulA(0.5)
	for _, l := range l.sources {
		// check if line segment is cw or ccw relative to light source, by calculating a signed area
		area := ((p2.X-p1.X)*(l.GetPos().Y-p1.Y) - (l.GetPos().X-p1.X)*(p2.Y-p1.Y))
		if area < 0 {
			col.Incr(l.ComputeColor(p, n, ownCol))
		}
	}

	return col.clamp()
}

func (l *Lighting) Draw() {
	for _, l := range l.sources {
		l.Draw()
	}
}
func (ol *OmniLight) SetColor(col Color) { ol.Col = col }
func (ol OmniLight) GetPos() V2          { return ol.Pos }
func (ol OmniLight) ComputeColor(at, normal V2, col Color) Color {
	col = col.Mul(ol.Col)
	v := ol.Pos.Sub(at)
	dist := v.Len()
	l := v.DivA(dist)

	diffuse := normal.NormDot(l)
	diffuse = math.Abs(diffuse)
	diffuse *= 200 * ol.Strength / (dist*dist + 0.001)

	if diffuse > 1 {
		diffuse = 1
	}

	c := col.MulA(diffuse)
	return c
}
func (ol OmniLight) Draw() {
	rl.DrawCircleGradient(int32(ol.Pos.X), int32(ol.Pos.Y), float32(ol.Strength), ColorToRlColorFade(ol.Col, 0.4), rl.NewColor(0, 0, 0, 0))
}
func (sl *SpotLight) ComputeColor(at, normal V2, col Color) Color {
	dir2light := at.Sub(sl.Pos).Norm()
	liDir := v.RotV(sl.Dir)
	cone := math.Cos(sl.Angle * rl.Deg2rad / 2)
	dist2 := at.Sub(sl.Pos).Len2()
	if dir2light.NormDot(liDir) > cone && dist2 < sl.Strength*sl.Strength {
		return sl.OmniLight.ComputeColor(at, normal, col).MulA(0.8)
	} else {
		return newColorRGB(0, 0, 0)
	}
}
func (sl *SpotLight) Draw() {
	pos := make([]float32, 2)
	pos[0], pos[1] = float32(sl.Pos.X), float32(Game.gH-sl.Pos.Y)
	size := make([]float32, 1)
	size[0] = float32(sl.Strength)
	p1 := sl.Pos.Add(v.RotV(sl.Dir + sl.Angle/2).MulA(sl.Strength))
	p2 := sl.Pos.Add(v.RotV(sl.Dir - sl.Angle/2).MulA(sl.Strength))

	rl.SetShaderValue(slshader, rl.GetShaderLocation(slshader, "pos"), pos, rl.ShaderUniformVec2)
	rl.SetShaderValue(slshader, rl.GetShaderLocation(slshader, "size"), size, rl.ShaderUniformFloat)
	rl.BeginShaderMode(slshader)
	_triangle(sl.Pos, p1, p2, ColorToRlColor(sl.Col))
	//rl.DrawRectangleV(rl.Vector2{0, 0}, rl.Vector2{float32(Game.gW), float32(Game.gH)}, rl.Yellow)

	rl.EndShaderMode()
}
