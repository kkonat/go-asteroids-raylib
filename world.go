package main

import (
	"fmt"

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
func (w *world) setupCamera(position, target, up rl.Vector3) {
	w.camera = rl.NewCamera3D(position, target, up, 10.0, rl.CameraProjection(rl.CameraOrthographic))
	rl.SetCameraMode(w.camera, rl.CameraFree)

	w.camMat = rl.GetCameraMatrix(w.camera)
}
func (w *world) drawBackground() {
	const size = 5.0

	for i := float32(-size); i < float32(size+1); i += 1.0 {
		rl.DrawLine3D(rl.NewVector3(i, -size, 0), rl.NewVector3(i, size, 0), rl.DarkGray)
		rl.DrawLine3D(rl.NewVector3(-size, i, 0), rl.NewVector3(size, i, 0), rl.DarkGray)
	}
}

func (w *world) drawObjects() {
	for _, o := range w.objects {
		m := o.getTransfMatrix()
		mv := rl.MatrixMultiply(m, w.camMat)
		rl.SetMatrixModelview(mv)
		fmt.Println(mv)
		o.Draw()
	}
	rl.SetMatrixModelview(w.camMat)
	fmt.Println(w.camMat)
}

func (w *world) animate(time_inc float32) {
	for _, o := range w.objects {
		o.Move(time_inc)
	}
}
