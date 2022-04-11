package soundmanager

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO
// for playback/panning https://bitbucket.org/StephenPatrick/go-winaudio/src/master/winaudio/
//
// for ogg https://github.com/mccoyst/ogg
// for mp3 https://github.com/hajimehoshi/go-mp3

type sound struct {
	rlSound   rl.Sound
	vol       float32
	maxVol    float32
	fadeCount int
	playtime  int
}

type SoundManager struct {
	sounds     []*sound
	fade       bool
	Mute       bool
	loopSounds []int
}
type SoundFile struct {
	Fname string
	Vol   float32
	Pitch float32
}

func NewSoundManager(mute bool, SoundFiles map[int]SoundFile) *SoundManager {
	rl.InitAudioDevice()
	sm := new(SoundManager)
	sm.sounds = make([]*sound, len(SoundFiles))

	for i, sf := range SoundFiles {
		snd := new(sound)
		rlSnd := rl.LoadSound(sf.Fname)
		if rlSnd.Stream.Buffer == nil {
			panic("can't load sound, Oh no! " + sf.Fname)
		}
		snd.rlSound = rlSnd
		snd.maxVol, snd.vol = sf.Vol, sf.Vol
		rl.SetSoundPitch(rlSnd, sf.Pitch)
		rl.SetSoundVolume(rlSnd, sf.Vol)
		sm.sounds[i] = snd
	}

	sm.Mute = mute
	return sm
}

func (sm *SoundManager) EnableLoops(sounds ...int) {
	sm.loopSounds = append(sm.loopSounds, sounds...)
}
func (sm *SoundManager) Update() {
	for s := range sm.loopSounds {
		if !sm.IsPlaying(s) {
			sm.Play(s)
		}
	}

	sm.DoFade() // fade out sounds if needed
}
func (sm *SoundManager) Stop(idx int) {

	// As it turns out, in raylib you can not simply stop playing sound, because you will hear a pop
	// You have to gently fade the sound out and when its barely audible, you can stop it

	// So, it was like:
	//rl.StopSound(sm.sounds[idx])
	//fmt.Println("fade:", idx)

	// but now, it has to be:
	sm.fade = true                // tell the doFade() that ther's something to fade out
	sm.sounds[idx].fadeCount = 60 // fade out in this many steps

}
func (sm *SoundManager) StopAll() {
	for i := range sm.sounds {
		sm.sounds[i].fadeCount = 60
	}
	sm.fade = true
}

// so this function is called every frame to check if there's anything to fade out
func (sm *SoundManager) DoFade() {
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
func (sm *SoundManager) Play(idx int) {
	if !sm.Mute {
		sm.sounds[idx].fadeCount = 0
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		sm.sounds[idx].vol = sm.sounds[idx].maxVol
		rl.PlaySound(sm.sounds[idx].rlSound)

	}
}
func (sm *SoundManager) PlayFor(idx, cycles int) {
	if !sm.Mute {
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
func (sm *SoundManager) PlayM(idx int) {
	if !sm.Mute {
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		rl.PlaySoundMulti(sm.sounds[idx].rlSound)
	}
}
func (sm *SoundManager) PlayPM(idx int, pitch float32) {
	if !sm.Mute {
		rl.SetSoundVolume(sm.sounds[idx].rlSound, sm.sounds[idx].maxVol)
		rl.SetSoundPitch(sm.sounds[idx].rlSound, pitch)
		rl.PlaySoundMulti(sm.sounds[idx].rlSound)
	}
}
func (sm *SoundManager) IsPlaying(idx int) bool {
	return rl.IsSoundPlaying(sm.sounds[idx].rlSound)
}
func (sm *SoundManager) UnloadAll() {
	for _, s := range sm.sounds {
		rl.StopSound(s.rlSound)
		rl.UnloadSound(s.rlSound)
	}
}
