package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	sprAlien0 = iota
)

type sprite struct {
	tex    rl.Texture2D
	w, h   int32
	frames uint8
	cx, cy int32
	shader rl.Shader
}
type spriteManager struct {
	sprites []*sprite
}

func newSprite(tex rl.Texture2D, w, h, cx, cy int32, frames uint8) *sprite {
	spr := new(sprite)
	spr.tex = tex
	spr.w, spr.h = w, h
	spr.cx, spr.cy = cx, cy
	spr.frames = frames
	return spr
}
func newSpriteManager() *spriteManager {

	sm := new(spriteManager)
	sm.loadTexture("res/ufo3.png", 32, 32, 16, 16, 8)
//	sm.sprites[sprAlien0].shader = rl.LoadShader("shaders/base.vs", "shaders/blur.fs")

	//colDiffuse
//	sh := sm.sprites[sprAlien0].shader
//	col := make([]float32, 4)
//	rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "blur"), col, rl.ShaderUniformIvec4)

	return sm
}
func (sm *spriteManager) loadTexture(filename string, w, h, cx, cy int32, frames uint8) int {

	var r int
	tex := rl.LoadTexture(filename)
	s := newSprite(tex, w, h, cx, cy, frames)
	sm.sprites = append(sm.sprites, s)
	return r
}

func (sm *spriteManager) drawSprite(id int, frame uint8, x, y int, scale, rot float64) {
	src := rl.NewRectangle(float32(int32(frame)*sm.sprites[id].w), 0,
		float32(sm.sprites[id].w), float32(sm.sprites[id].h))
	targ := rl.NewRectangle(float32(x), float32(y),
		float32(float64(sm.sprites[id].w)*scale),
		float32(float64(sm.sprites[id].h)*scale))

	//	rl.BeginShaderMode(sm.sprites[id].shader)

	rl.DrawTexturePro(
		sm.sprites[id].tex, // from texture
		src,                // crop source
		targ,               // screen dest
		rl.Vector2{X: float32(sm.sprites[id].cx), Y: float32(sm.sprites[id].cy)}, // offset
		float32(rot), // rotation
		rl.White)

	//	rl.EndShaderMode()
}

func (sm *spriteManager) unloadAll() {

	for _, s := range sm.sprites {
		rl.UnloadTexture(s.tex)
	}

}
