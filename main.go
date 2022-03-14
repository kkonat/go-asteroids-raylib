package main

import "C"
import (
	rl "github.com/gen2brain/raylib-go/raylib"
	// ns "github.com/ojrac/opensimplex-go"
)

type object interface {
	Draw()
	Move(time_inc float32)
	getTransfMatrix() rl.Matrix
}
type Cube struct {
	Model  rl.Model
	time   float32
	pos    rl.Vector2
	speed  float32
	color  rl.Color
	trMtrx rl.Matrix
}

func NewCube(pos rl.Vector2, size, speed float32, color rl.Color) *Cube {
	c := new(Cube)
	cube := rl.GenMeshCube(size, size, size)
	c.Model = rl.LoadModelFromMesh(cube)
	c.pos = pos
	c.speed = speed
	c.time = 0.0
	c.color = color
	return c
}
func (c *Cube) Draw() {
	//rl.DrawModelEx(c.Model, rl.NewVector3(0, 0, 0), rl.NewVector3(0, 0, 1), c.time, rl.NewVector3(1, 1, 1), c.color*0.5)
	//rl.DrawModelWiresEx(c.Model, rl.NewVector3(c.pos.X, c.pos.Y, 0), rl.NewVector3(0, 0, 1), c.time, rl.NewVector3(1, 1, 1), c.color)
	rl.DrawModelWires(c.Model, rl.NewVector3(c.pos.X, c.pos.Y, 0), 1.0, c.color)
}
func (c *Cube) Move(time_inc float32) {
	c.time += time_inc * c.speed
	c.trMtrx = rl.MatrixRotateZ(c.time / 180 * 3.14)
}

func (c *Cube) getTransfMatrix() rl.Matrix { return c.trMtrx }

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
