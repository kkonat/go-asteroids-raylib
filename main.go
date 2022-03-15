package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	game := newGame(1440, 720)

	for !rl.WindowShouldClose() {
		if !game.sm.isPlaying(0) {
			game.sm.play(game.sm.sSpace)

			fmt.Println("started")
		}

		if rl.IsKeyPressed('A') {
			game.sm.playM(game.sm.sOinx)
		}
		if rl.IsKeyPressed('S') {
			if !game.sm.mute {
				game.sm.stop(0)
			}
			game.sm.mute = !game.sm.mute

		}

		game.drawGame()
	}
	game.finalize()

}
