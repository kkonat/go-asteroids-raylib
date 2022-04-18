package main

import (
	"fmt"
	qt "rlbb/lib/quadtree"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	debugMemstats         = true
	debugRocksQt          = false
	debugShipPos          = false
	degubDrawMissileLines = false
)

// -- debug
func DrawQt(qt *qt.QuadTree[*Rock]) {
	if len(qt.Objects) != 0 {
		rl.DrawRectangleLines(qt.Bounds.X+2, qt.Bounds.Y+2, qt.Bounds.W-4, qt.Bounds.H-4, rl.DarkGray)
		str := fmt.Sprintf("#%d", len(qt.Objects))
		rl.DrawText(str, qt.Bounds.X+2, qt.Bounds.Y+20, 16, rl.Gray)
	}
	for i := 0; i < 4; i++ {
		if qt.Nodes[i] != nil {
			DrawQt(qt.Nodes[i])
		}
	}
}
func (gme *game) debugQt() {
	if debug {
		if debugShipPos {
			str := fmt.Sprintf("[%d,%d]", int32(gme.ship.pos.X), int32(gme.ship.pos.Y))
			rl.DrawText(str, int32(gme.ship.pos.X), int32(gme.ship.pos.Y), 20, rl.Gray)
		}

		if debugRocksQt {

			potCols := gme.RocksQt.MayCollide(gme.ship.shape.bRect, minidist2)

			for _, c := range potCols {
				rl.DrawRectangleLines(c.BRect().X, c.BRect().Y, c.BRect().W, c.BRect().H, rl.DarkBrown)
			}
		}

		var line int32 = 16
		inc := func(l *int32) int32 { *l += 16; return *l }

		//printMemoryUsage
		if debugMemstats {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			str := fmt.Sprintf("Alloc = %v MiB", m.Alloc/1024/1024)
			rl.DrawText(str, 0, inc(&line), 16, rl.Gray)
			str = fmt.Sprintf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
			rl.DrawText(str, 0, inc(&line), 16, rl.Gray)
			str = fmt.Sprintf("\tSys = %v MiB", m.Sys/1024/1024)
			rl.DrawText(str, 0, inc(&line), 16, rl.Gray)
			str = fmt.Sprintf("\tNumGC = %v\n", m.NumGC)
			rl.DrawText(str, 0, inc(&line), 16, rl.Gray)
		}
		rl.DrawTextEx(gme.vectorFont, "DEBUG MODE", rl.Vector2{X: float32(720 - rl.MeasureText("DEBUG", 99)), Y: float32(590)}, 99, 0, rl.DarkPurple)

	}
}
