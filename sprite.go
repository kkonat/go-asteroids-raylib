package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	sprAlien0 = iota
	ammoCrate
)

type sprite struct {
	tex    rl.Texture2D
	tint   rl.Color
	w, h   int32
	frames uint8
	cx, cy int32
	shader rl.Shader
}
type spriteManager struct {
	sprites map[int]*sprite
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
	sm.sprites = make(map[int]*sprite)
	sm.loadSprite(sprAlien0, "res/ufo3.png", 32, 32, 16, 16, 8)
	sm.setTint(sprAlien0, rl.ColorFromHSV(193, 58, 162))
	sm.loadSprite(ammoCrate, "res/ammocrate.png", 32, 32, 16, 16, 1)
	sm.setTint(ammoCrate, rl.ColorFromHSV(33, 58, 122))


	//	sm.sprites[sprAlien0].shader = rl.LoadShader("shaders/base.vs", "shaders/blur.fs")
	//	sh := sm.sprites[sprAlien0].shader
	//	col := make([]float32, 4)
	//	rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "blur"), col, rl.ShaderUniformIvec4)

	return sm
}
func (sm *spriteManager) setTint(idx int, col rl.Color) {
	sm.sprites[idx].tint = col
}
func (sm *spriteManager) loadSprite(idx int, filename string, w, h, cx, cy int32, frames uint8) {
	tex := rl.LoadTexture(filename)
	if tex.ID == 0 {
		msg := fmt.Sprintf("can't load []%s]", filename)
		panic(msg)
	}
	s := newSprite(tex, w, h, cx, cy, frames)

	sm.sprites[idx] = s
	sm.sprites[idx].tint = rl.White
}

func (sm *spriteManager) drawSprite(idx int, frame uint8, x, y int, scale, rot float64) {
	src := rl.NewRectangle(float32(int32(frame)*sm.sprites[idx].w), 0,
		float32(sm.sprites[idx].w), float32(sm.sprites[idx].h))
	targ := rl.NewRectangle(float32(x), float32(y),
		float32(float64(sm.sprites[idx].w)*scale),
		float32(float64(sm.sprites[idx].h)*scale))

	//	rl.BeginShaderMode(sm.sprites[id].shader)

	rl.DrawTexturePro(
		sm.sprites[idx].tex, // from texture
		src,                 // crop source
		targ,                // screen dest
		rl.Vector2{X: float32(sm.sprites[idx].cx), Y: float32(sm.sprites[idx].cy)}, // offset
		float32(rot), // rotation
		sm.sprites[idx].tint)

	//	rl.EndShaderMode()
}

func (sm *spriteManager) unloadAll() {

	for _, s := range sm.sprites {
		rl.UnloadTexture(s.tex)
	}

}
