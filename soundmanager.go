package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO
// for playback/panning https://bitbucket.org/StephenPatrick/go-winaudio/src/master/winaudio/
//
// for ogg https://github.com/mccoyst/ogg
// for mp3 https://github.com/hajimehoshi/go-mp3

const (
	sSpace = iota
	sScore
	sMissilesDlvrd
	sThrust
	sExpl
	sLaunch
	sShieldsLow
	sAmmoLow
	sOinx
	sExplodeShip
	sScratch
	sChargeUp
	sForceField
)

var soundFiles = map[int]struct {
	fname      string
	vol, pitch float32
}{
	// Id			  filename        vol  pitch
	sSpace:         {"res/space.ogg", 0.5, 1.0},
	sScore:         {"res/score.mp3", 0.1, 1.0},
	sMissilesDlvrd: {"res/missiles-delivered.ogg", 0.5, 1.0},
	sThrust:        {"res/thrust.ogg", 0.5, 1.0},
	sExpl:          {"res/expl.ogg", 0.5, 0.65},
	sLaunch:        {"res/launch.ogg", 0.5, 1.0},
	sShieldsLow:    {"res/warning-shields-low.ogg", 0.3, 1.0},
	sAmmoLow:       {"res/warning-ammo-low.ogg", 0.3, 1.0},	
	sOinx:          {"res/oinxL.ogg", 0.5, 1.0},
	sExplodeShip:   {"res/shipexplode.ogg", 1.0, 1.0},
	sScratch:       {"res/metalScratch.ogg", 0.2, 1.0},
	sChargeUp:      {"res/chargeup.ogg", 0.2, 1.0},
	sForceField:    {"res/forcefield2.ogg", 0.5, 1.0},
}

type sound struct {
	rlSound   rl.Sound
	vol       float32
	maxVol    float32
	fadeCount int
	playtime  int
}

type soundManager struct {
	sounds []*sound
	fade   bool
	mute   bool
}

func newSoundManager(mute bool) *soundManager {
	rl.InitAudioDevice()
	sm := new(soundManager)
	sm.sounds = make([]*sound, len(soundFiles))

	for i, sf := range soundFiles {
		snd := new(sound)
		rlSnd := rl.LoadSound(sf.fname)
		if rlSnd.Stream.Buffer == nil {
			panic("can't load sound, Oh no! " + sf.fname)
		}
		snd.rlSound = rlSnd
		snd.maxVol, snd.vol = sf.vol, sf.vol
		rl.SetSoundPitch(rlSnd, sf.pitch)
		rl.SetSoundVolume(rlSnd, sf.vol)
		sm.sounds[i] = snd
	}

	sm.mute = mute
	return sm
}

func (sm *soundManager) stop(idx int) {

	// As it turns out, in raylib you can not simply stop playing sound, because you will hear a pop
	// You have to gently fade the sound out and when its barely audible, you can stop it

	// So, it was like:
	//rl.StopSound(sm.sounds[idx])
	//fmt.Println("fade:", idx)

	// but now, it has to be:
	sm.fade = true                // tell the doFade() that ther's something to fade out
	sm.sounds[idx].fadeCount = 60 // fade out in this many steps

}
func (sm *soundManager) stopAll() {
	for i := range sm.sounds {
		sm.sounds[i].fadeCount = 60
	}
	sm.fade = true
}

// so this function is called every frame to check if there's anything to fade out
func (sm *soundManager) doFade() {
	var notFading = 0
	if sm.fade {
		for i, s := range sm.sounds {
			c := s.fadeCount
			if c > 0 {
				sm.sounds[i].vol *= 0.95
				rl.SetSoundVolume(sm.sounds[i].rlSound, sm.sounds[i].vol)
				c--
				sm.sounds[i].fadeCount = c
				if c == 0 {
					notFading++
					rl.StopSound(sm.sounds[i].rlSound)
				}
			} else {
				notFading++
			}
		}
		if notFading == len(sm.sounds) {
			sm.fade = false
		}
	}
}
func (sm *soundManager) play(idx int) {
	if !sm.mute {
		sm.sounds[idx].fadeCount = 0
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		sm.sounds[idx].vol = sm.sounds[idx].maxVol
		rl.PlaySound(sm.sounds[idx].rlSound)

	}
}
func (sm *soundManager) playFor(idx, cycles int) {
	if !sm.mute {
		if sm.sounds[idx].playtime == 0 {
			sm.sounds[idx].fadeCount = 0
			sm.sounds[idx].playtime = cycles
			rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
			sm.sounds[idx].vol = sm.sounds[idx].maxVol

			rl.PlaySound(sm.sounds[idx].rlSound)
		} else {
			sm.sounds[idx].playtime--
		}

	}
}
func (sm *soundManager) playM(idx int) {
	if !sm.mute {
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		rl.PlaySoundMulti(sm.sounds[idx].rlSound)
	}
}
func (sm *soundManager) playPM(idx int, pitch float32) {
	if !sm.mute {
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		rl.SetSoundPitch(sm.sounds[idx].rlSound, pitch)
		rl.PlaySoundMulti(sm.sounds[idx].rlSound)
	}
}
func (sm *soundManager) isPlaying(idx int) bool {
	return rl.IsSoundPlaying(sm.sounds[idx].rlSound)
}
func (sm *soundManager) unloadAll() {
	for _, s := range sm.sounds {
		rl.StopSound(s.rlSound)
		rl.UnloadSound(s.rlSound)
	}
}
