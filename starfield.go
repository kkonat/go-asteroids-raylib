package main

import (
	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed res/Space.png
var space []byte

type star struct {
	y        int32
	x, speed float32
	r        uint8
}

type starfield struct {
	stars       []star
	w, h        int32
	starfTex    rl.Texture2D
	shader      rl.Shader
	time        []float32 // shader uniforms
	iResolution []float32 // shader uniforms
	timeLoc     int32
}

func newStarfield(w, h int32, time []float32) *starfield {
	sf := new(starfield)
	sf.time = time
	sf.w, sf.h = w, h

	// img := rl.LoadImage("res/Space.png")
	img := rl.LoadImageFromMemory(".png", space, int32(len(space)))

	rl.ImageCrop(img, rl.Rectangle{X: 0, Y: 0, Width: float32(w), Height: float32(h)})
	sf.starfTex = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	// sf.shader = rl.LoadShader("shaders/base.vs", "shaders/starfield.fs")
	sf.shader = load_shaders("shaders/base.vs", "shaders/starfield.fs")

	sf.time[0] = 9321
	sf.iResolution = make([]float32, 2)
	sf.iResolution[0], sf.iResolution[1] = float32(sf.w), float32(sf.h)
	sf.timeLoc = rl.GetShaderLocation(sf.shader, "time")

	rl.SetShaderValue(sf.shader, rl.GetShaderLocation(sf.shader, "iResolution"), sf.iResolution, rl.ShaderUniformVec2)
	return sf
}
func (sf *starfield) draw() {

	rl.SetShaderValue(sf.shader, sf.timeLoc, sf.time, rl.ShaderUniformFloat)

	rl.BeginShaderMode(sf.shader)
	rl.DrawTexture(sf.starfTex, 0, 0, rl.White)
	rl.EndShaderMode()

	//	rl.DrawCircleGradient(1655, 400, 900, rl.NewColor(30, 20, 105, 255), rl.NewColor(0, 0, 0, 0))
	rl.DrawCircleLines(1655, 400, 400, rl.NewColor(40, 10, 140, 255))
	rl.DrawCircle(1655, 400, 400, rl.NewColor(20, 0, 70, 255))

}
