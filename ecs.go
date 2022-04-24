package main

type Motion struct {
	Pos, Speed    V2
	Rot, RotSpeed float64
}

func (m *Motion) Move(dt float64) {
	dv := m.Speed.MulA(dt)
	m.Pos.Incr(dv)
	m.Rot += m.RotSpeed * dt
}

type movable interface {
	Move(dt float64)
}
type drawable interface {
	Draw()
}
type collidable interface {
	BRect()
}

type Object struct{}

type Static struct {
	Object
}

type Mobile struct {
	Object
	Motion
}

func (m *Mobile) Draw() {}
