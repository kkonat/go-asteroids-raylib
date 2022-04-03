package main

import (
	"fmt"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	debugRockcount        = true
	debugMemstats         = true
	debugRocksQt          = false
	debugShipPos          = true
	degubDrawMissileLines = false
)

// -- debug
func drawQt(qt *QuadTree[RockListEl]) {
	if len(qt.Objects) != 0 {
		rl.DrawRectangleLines(qt.Bounds.x+2, qt.Bounds.y+2, qt.Bounds.w-4, qt.Bounds.h-4, rl.DarkGray)
		str := fmt.Sprintf("#%d", len(qt.Objects))
		rl.DrawText(str, qt.Bounds.x+2, qt.Bounds.y+20, 16, rl.Gray)
	}
	for i := 0; i < 4; i++ {
		if qt.Nodes[i] != nil {
			drawQt(qt.Nodes[i])
		}
	}
}
func (gme *game) debugQt() {
	if debug {
		if debugShipPos {
			// str := fmt.Sprintf("[%d,%d]",int32(gme.ship.pos.x),int32(gme.ship.pos.y))
			// rl.DrawText(str, int32(gme.ship.pos.x),int32(gme.ship.pos.y), 20, rl.Gray)
		}
		//x, y := int32(gme.ship.pos.x), int32(gme.ship.pos.y)
		//	shipRect := Rect{x, y, 20, 20}

		// potCols := gme.RocksQt.MayCollide(shipRect)
		// for _, c := range potCols {
		// 	rl.DrawRectangleLines(c.bRect().x, c.bRect().y, c.bRect().w, c.bRect().h, rl.Beige)
		// }

		var line int32 = 16
		inc := func(l *int32) int32 { *l += 16; return *l }

		if debugRockcount {
			str := fmt.Sprintf("rocks len = %v", gme.rocks.Len)
			rl.DrawText(str, 0, inc(&line), 16, rl.White)
		}

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
		if debugRocksQt {
			drawQt(gme.RocksQt)
			// for i, missile := range gme.missiles {
			// 	largerCircle = newCircleV2(missile.pos, 10)
			// 	potCols = gme.qt.MayCollide(largerCircle)
			// 	for _, c := range potCols {
			// 		rl.DrawRectangleLines(c.rect.x, c.rect.y, c.rect.w, c.rect.h, rl.DarkGreen)
			// 		str := fmt.Sprintf("[%d]", i)
			// 		rl.DrawText(str, c.rect.x+int32(i*16), c.rect.y, 16, rl.Lime)
			// 	}
			// }
		}
		rl.DrawTextEx(vectorFont, "DEBUG MODE", rl.Vector2{X: float32(720 - rl.MeasureText("DEBUG", 99)), Y: float32(590)}, 99, 0, rl.DarkPurple)

	}
}
