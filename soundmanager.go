package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type soundManager struct {
	sSpace, sOinx int
	sounds        []rl.Sound
	mute          bool
}

func newSoundManager() *soundManager {
	rl.InitAudioDevice()
	sm := new(soundManager)
	sm.sSpace = sm.loadSound("res/space.ogg", 1.0, 0.52)
	sm.sOinx = sm.loadSound("res/oinx.wav", 1.0, 1.0)
	sm.mute = false
	return sm
}

func (sm *soundManager) loadSound(fname string, volume, pitch float32) int {
	snd := rl.LoadSound(fname)

	rl.SetSoundPitch(snd, pitch)
	rl.SetSoundVolume(snd, volume)
	sm.sounds = append(sm.sounds, snd)

	return len(sm.sounds) - 1
}
func (sm *soundManager) stop(idx int) {
	rl.StopSound(sm.sounds[idx])
}
func (sm *soundManager) play(idx int) {
	if !sm.mute {
		rl.PlaySound(sm.sounds[idx])
	}
}
func (sm *soundManager) playM(idx int) {
	if !sm.mute {
		rl.PlaySoundMulti(sm.sounds[idx])
	}
}
func (sm *soundManager) isPlaying(idx int) bool {
	return rl.IsSoundPlaying(sm.sounds[idx])
}
func (sm *soundManager) unloadAll() {
	for _, s := range sm.sounds {
		rl.StopSound(s)
		rl.UnloadSound(s)
	}
}
