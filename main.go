package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// By default, Go programs run with GOMAXPROCS set to the number
	// of cores available; in prior releases it defaulted to 1.
	// Starting from Go 1.5, the default value is the number of cores.
	// You only need to explicitly set it if you are not okay with this
	// in newer Go versions.
	//runtime.GOMAXPROCS(8)

	rl.SetTraceLog(rl.LogAll)
	// GetWindowHandle();
	game := newGame(1440, 720)

	for !rl.WindowShouldClose() {

		game.processKeys()
		game.processMouse()
		game.playMessages()
		game.drawAndUpdate()

	}
	game.finalize()
}




