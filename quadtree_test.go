package main

import "testing"

func TestQuadtreeInsert(t *testing.T) {
	r := Rect{0, 0, 320, 320}
	qt := NewQuadTree(0, r)

	if qt.Level != 0 {
		t.Error("wrong level")
	}

	c := newCircleV2(V2{100, 100}, 50)
	c1 := newCircleV2(V2{110, 100}, 10)
	qt.Insert(c)
	cw := qt.MayCollide(c1)
	if c != cw[0] {
		t.Error("different objects")
	}
	qt.Insert(newCircleV2(V2{100, 110}, 10))

	if len(qt.MayCollide(c1)) != 2 {
		t.Error("not enough collissions")
	}
	for i := 0; i < 14; i++ {
		qt.Insert(newCircleV2(V2{290 + float64(i), 300 - float64(i)}, 1))
	}
	if len(qt.MayCollide(c1)) != 2 {
		t.Error("too many collissions")
	}

	qt.Clear()
	if qt.Nodes[0] != nil || qt.Nodes[1] != nil ||
		qt.Nodes[2] != nil || qt.Nodes[3] != nil {
		t.Error("not cleared")
	}
}
