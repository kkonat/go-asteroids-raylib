package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type soundManager struct {
	sSpace, sOinx, sThrust int
	sExpl, sLaunch         int
	sounds                 []rl.Sound
	volumes                []float32
	maxvolumes             []float32
	fadecount              []int32
	fade                   bool
	mute                   bool
}

func newSoundManager(mute bool) *soundManager {
	rl.InitAudioDevice()
	sm := new(soundManager)
	sm.sSpace = sm.loadSound("res/space.ogg", 0.5, 0.52)
	sm.sOinx = sm.loadSound("res/oinxL.ogg", 0.7, 1.0)
	sm.sThrust = sm.loadSound("res/thrust.ogg", 1.0, 1.0)
	sm.sExpl = sm.loadSound("res/expl.ogg", 0.5, 0.65)
	sm.sLaunch = sm.loadSound("res/launch.ogg", 0.5, 1.0)
	sm.mute = mute
	return sm
}

func (sm *soundManager) loadSound(fname string, volume, pitch float32) int {
	snd := rl.LoadSound(fname)

	rl.SetSoundPitch(snd, pitch)
	rl.SetSoundVolume(snd, volume)
	sm.sounds = append(sm.sounds, snd)
	sm.volumes = append(sm.volumes, volume)
	sm.maxvolumes = append(sm.maxvolumes, volume)
	sm.fadecount = append(sm.fadecount, 0)
	return len(sm.sounds) - 1
}
func (sm *soundManager) stop(idx int) {

	// As it turns out, in raylib you can not simply stop playing sound, because you will hear a pop
	// You have to gently fade the sound out and when its barely audible, you can stop it

	// So, it was like:
	//rl.StopSound(sm.sounds[idx])
	//fmt.Println("fade:", idx)

	// but now, it has to be:
	sm.fade = true         // tell the doFade() that ther's something to fade out
	sm.fadecount[idx] = 60 // fade out in this many steps

}

// so this function is called every frame to check if there's anything to fade out
func (sm *soundManager) doFade() {
	var notFading = 0
	if sm.fade {
		for i, c := range sm.fadecount {
			if c > 0 {

				v := sm.volumes[i] * 0.95
				sm.volumes[i] = v
				rl.SetSoundVolume(sm.sounds[i], v)
				c--
				sm.fadecount[i] = c
				if c == 0 {
					notFading++
					rl.StopSound(sm.sounds[i])
				}
			} else {
				notFading++
			}
		}
		if notFading == len(sm.fadecount) {
			sm.fade = false
		}
	}
}
func (sm *soundManager) play(idx int) {
	if !sm.mute {
		sm.fadecount[idx] = 0
		rl.SetSoundVolume(sm.sounds[idx], sm.maxvolumes[idx])
		sm.volumes[idx] = sm.maxvolumes[idx]

		rl.PlaySound(sm.sounds[idx])

	}
}
func (sm *soundManager) playM(idx int) {
	if !sm.mute {
		rl.PlaySoundMulti(sm.sounds[idx])
	}
}
func (sm *soundManager) playPM(idx int, pitch float32) {
	if !sm.mute {
		rl.SetSoundPitch(sm.sounds[idx], pitch)
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
