package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type world struct {
	camera rl.Camera3D
	camMat rl.Matrix

	objects []object
}

func newWorld() *world {
	return new(world)
}
func (w *world) AddObject(o object) {
	w.objects = append(w.objects, o)
}

func (w *world) drawObjects() {
	for _, o := range w.objects {
		o.Draw()
	}
	rl.SetMatrixModelview(w.camMat)

}

func (w *world) animate(time_inc float32) {
	for _, o := range w.objects {
		o.Move(time_inc)
	}
}
