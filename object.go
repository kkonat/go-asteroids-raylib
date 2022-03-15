package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type object interface {
	Draw()
	Move(time_inc float32)
}

type movable struct {
	pos, pivot rl.Vector2
	speed, rot float32
	trMtrx     rl.Matrix
}
