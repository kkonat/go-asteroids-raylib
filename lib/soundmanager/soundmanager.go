package soundmanager

import (
	"embed"
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed sounds/*
var soundsFS embed.FS

func load_sound(filename string) rl.Sound {
	ext := filename[len(filename)-4:]
	soundBytes, err := soundsFS.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read sound file: %v", err)
	}
	soundData := []byte(soundBytes)
	wave := rl.LoadWaveFromMemory(ext, soundData, int32(len(soundData)))
	sound := rl.LoadSoundFromWave(wave)

	return sound
}

// TODO
// for playback/panning https://bitbucket.org/StephenPatrick/go-winaudio/src/master/winaudio/
//
// for ogg https://github.com/mccoyst/ogg
// for mp3 https://github.com/hajimehoshi/go-mp3

// sound data structure
type sound struct {
	rlSound   rl.Sound // raylib sound id
	vol       float32  // current volume
	maxVol    float32  // maximum volume
	fadeCount int      // counter for fading out
	playtime  int      // play duration
}

// a data structure for cyclic voice messages
type VoiceMsg = struct {
	Lastplayed, // time (in seconds) when the message was last played
	Repeat int64 // repeat every n seconds
}

type SoundManager struct {
	sounds     []*sound         // a slice holding all loopSounds
	fade       bool             // global fade flag, if any sound needs fading this is set to true
	Mute       bool             // all sounds are muted
	loopSounds []int            // slice of ids of sounds, which need to be looped
	Msgs       map[int]VoiceMsg // map of sounds, which are used as cyclic voice messages
}

// used for sound manager initialization, holdes filenames, default volumes anf pitches of sounds to be loaded into the sound manager
type SoundFile struct {
	Fname string  // filename
	Vol   float32 // default volume
	Pitch float32 // default pitch
}

// Creeates a new sound manager
// mute determines the initial sound state and SoundFiles lists all filenames that will be used
func NewSoundManager(mute bool, SoundFiles map[int]SoundFile) *SoundManager {
	rl.InitAudioDevice()
	sm := new(SoundManager)
	sm.sounds = make([]*sound, len(SoundFiles))

	for i, sf := range SoundFiles {
		rlSnd := load_sound(sf.Fname)
		snd := new(sound)
		snd.rlSound = rlSnd
		snd.maxVol, snd.vol = sf.Vol, sf.Vol
		rl.SetSoundPitch(rlSnd, sf.Pitch)
		rl.SetSoundVolume(rlSnd, sf.Vol)
		sm.sounds[i] = snd
	}

	sm.Mute = mute
	return sm
}

// creates a list of sound, which will be played looped, like background music
func (sm *SoundManager) EnableLoops(sounds ...int) {
	sm.loopSounds = append(sm.loopSounds, sounds...)
}

// invoked on each frame in the game loop to chek which loop sounds are playing and which need to be restarted
func (sm *SoundManager) Update() {
	for s := range sm.loopSounds {
		if !rl.IsSoundPlaying(sm.sounds[s].rlSound) {
			sm.Play(s)
		}
	}

	sm.DoFade() // fade out sounds if needed
}

// stops a single sound
func (sm *SoundManager) Stop(idx int) {

	// As it turns out, in raylib you can not simply stop playing sound, because you will hear a pop
	// You have to gently fade the sound out and when its barely audible, you can stop it

	// So, it used to be like:
	//rl.StopSound(sm.sounds[idx])
	//fmt.Println("fade:", idx)

	// but now, it has to be:
	sm.fade = true                // tell the doFade() that ther's something to fade out
	sm.sounds[idx].fadeCount = 60 // fade out in this many steps

}

// stops all sounds
func (sm *SoundManager) StopAll() {
	for i := range sm.sounds {
		sm.sounds[i].fadeCount = 60
	}
	sm.fade = true
}

// called on every frame to check if there are any sounds that need to be faded out
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

// plays a single sound
func (sm *SoundManager) Play(sound int) {
	if !sm.Mute {
		sm.sounds[sound].fadeCount = 0
		rl.SetSoundVolume(sm.sounds[sound].rlSound, sm.sounds[sound].maxVol)
		sm.sounds[sound].vol = sm.sounds[sound].maxVol
		rl.PlaySound(sm.sounds[sound].rlSound)
	}
}

// plays a sound for a specific numver of frames
func (sm *SoundManager) PlayFor(sound, cycles int) {
	if !sm.Mute {
		if sm.sounds[sound].playtime == 0 {
			sm.sounds[sound].fadeCount = 0
			sm.sounds[sound].playtime = cycles

			rl.SetSoundVolume(sm.sounds[sound].rlSound, sm.sounds[sound].maxVol)
			sm.sounds[sound].vol = sm.sounds[sound].maxVol

			rl.PlaySound(sm.sounds[sound].rlSound)
		} else {
			sm.sounds[sound].playtime--
		}

	}
}

// plays a sound using multi-channel engine
func (sm *SoundManager) PlayM(sound int) {
	if !sm.Mute {
		rl.SetSoundVolume(sm.sounds[sound].rlSound, sm.sounds[sound].maxVol)
		rl.PlaySoundMulti(sm.sounds[sound].rlSound)
	}
}

// plays a sound using multi-channel engine with specific pitch
func (sm *SoundManager) PlayPM(sound int, pitch float32) {
	if !sm.Mute {
		rl.SetSoundVolume(sm.sounds[sound].rlSound, sm.sounds[sound].maxVol)
		rl.SetSoundPitch(sm.sounds[sound].rlSound, pitch)
		rl.PlaySoundMulti(sm.sounds[sound].rlSound)
	}
}

// cleanup function, unloads all sounds from memory
func (sm *SoundManager) UnloadAll() {
	for _, s := range sm.sounds {
		rl.StopSound(s.rlSound)
		rl.UnloadSound(s.rlSound)
	}
}

// plays cyclic message sounds, requires SoundManager.Msgs to be assigned first
func (sm *SoundManager) PlayCyclic(msg int) {
	t := time.Now().Local().Unix()
	if sm.Msgs[msg].Lastplayed == 0 || t-sm.Msgs[msg].Lastplayed > sm.Msgs[msg].Repeat {
		sm.PlayM(msg)
		sm.Msgs[msg] = VoiceMsg{t, sm.Msgs[msg].Repeat}
	}
}

// stops playing a cyclic message
func (sm *SoundManager) StopCyclic(msg int) {
	sm.Msgs[msg] = VoiceMsg{0, sm.Msgs[msg].Repeat}
}
