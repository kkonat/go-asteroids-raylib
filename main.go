package main

import "C"
import (
	rl "github.com/gen2brain/raylib-go/raylib"
	// ns "github.com/ojrac/opensimplex-go"
)

func main() {
	g := newGame(800, 600)

	w := newWorld()
	w.setupCamera(rl.NewVector3(0, 0, 5), rl.NewVector3(0, 0, 0), rl.NewVector3(0, 1, 0))

	w.AddObject(NewCube(rl.NewVector2(0, 0), 2.0, 1.0, rl.Blue))
	w.AddObject(NewCube(rl.NewVector2(3, 3), 3.0, -0.1, rl.Yellow))
	//w.AddObject(NewCube(rl.NewVector2(-3, -3), 1.0, -0, rl.Green))

	for !rl.WindowShouldClose() {
		g.drawGame(w)
		w.animate(1.0)
		// key := rl.GetKeyPressed()
		// if (key) != 0 {
		// 	fmt.Println("z=", z)
		// 	switch key {
		// 	case rl.KeyA:
		// 		z += 1
		// 	case rl.KeyS:
		// 		z = 0
		// 	case int32('D'):
		// 		z -= 1
		// 	}
		// }

	}
	g.finalize()
}
