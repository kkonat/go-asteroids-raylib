package main

import (
	"fmt"
)

const qtMaxObjects = 10
const qtMaxLevels = 5

type Rect struct {
	x, y int32
	w, h int32
}

type qtObj interface {
	comparable
	bRect() Rect
}

type QuadTree[T qtObj] struct {
	Level   int
	Objects []T
	Bounds  Rect
	Nodes   [4]*QuadTree[T]
	Total   int
}

func (qt *QuadTree[T]) TotalNodes() int {

	total := 0

	if qt.Nodes[0] != nil {
		for i := 0; i < 4; i++ {
			total += 1
			total += qt.Nodes[i].TotalNodes()
		}
	}

	return total

}
func (q *QuadTree[T]) Clear() {

	q.Objects = nil
	if q.Nodes[0] != nil {
		for i := 0; i < 4; i++ {
			q.Nodes[i].Clear()
			q.Nodes[i] = nil
		}
	}
	q.Total = 0

}
func newNode[T qtObj](pLevel int, pBounds Rect) *QuadTree[T] {
	n := new(QuadTree[T])
	n.Objects = make([]T, 0, qtMaxObjects)
	n.Level = pLevel
	n.Bounds = pBounds
	return n
}
func (q *QuadTree[T]) split() {
	subW, subH := q.Bounds.w/2, q.Bounds.h/2
	x, y := q.Bounds.x, q.Bounds.y

	q.Nodes[0] = newNode[T](q.Level+1, Rect{x, y, subW, subH})
	q.Nodes[1] = newNode[T](q.Level+1, Rect{x + subW, y, subW, subH})
	q.Nodes[2] = newNode[T](q.Level+1, Rect{x, y + subH, subW, subH})
	q.Nodes[3] = newNode[T](q.Level+1, Rect{x + subW, y + subH, subW, subH})
}

const (
	qTopLeft     = 0
	qTopRight    = 1
	qBottomLeft  = 2
	qBottomRight = 3
	qDoesntFit   = 5
)

func (q *QuadTree[T]) getQuadrant(r Rect) int {

	quadrant := qDoesntFit

	cx := q.Bounds.x + (q.Bounds.w / 2)
	cy := q.Bounds.y + (q.Bounds.h / 2)

	fitsInTop := (r.y+r.h < cy)
	fitsInBottom := (r.y > cy)
	// left quadrants
	if r.x+r.w < cx {
		if fitsInTop {
			quadrant = qTopLeft
		} else if fitsInBottom {
			quadrant = qBottomLeft
		}
	} else if r.x > cx { // right quadrants
		if fitsInTop {
			quadrant = qTopRight
		} else if fitsInBottom {
			quadrant = qBottomRight
		}
	}
	return quadrant
}
func (q *QuadTree[T]) find(obj T) bool {
	found := false
	if q.Nodes[0] == nil {
		for _, o := range q.Objects {
			if o == obj {
				return true
				//break
			}
		}
	} else {
		found = found || q.Nodes[0].find(obj)
		found = found || q.Nodes[1].find(obj)
		found = found || q.Nodes[2].find(obj)
		found = found || q.Nodes[3].find(obj)
	}
	return found
}
func (q *QuadTree[T]) Remove(objTbRemved T) bool {
	r := q.rem(objTbRemved, false)
	if !r {
		fmt.Printf(" >>QT el not removed")
		//q.rem(objTbRemved, true)
	}
	return r
}
func (q *QuadTree[T]) rem(objTbRemved T, debug bool) bool {
	removed := false
	if q.Nodes[0] == nil {
		for i, o := range q.Objects {
			if debug {
				fmt.Printf("_%p_ ", o)
			}
			if o == objTbRemved {
				q.Objects = append(q.Objects[:i], q.Objects[i+1:]...)
				fmt.Print(" X removed X ")
				q.Total--
				return true
				//break
			}
		}
	} else {
		removed = removed || q.Nodes[0].rem(objTbRemved, debug)
		removed = removed || q.Nodes[1].rem(objTbRemved, debug)
		removed = removed || q.Nodes[2].rem(objTbRemved, debug)
		removed = removed || q.Nodes[3].rem(objTbRemved, debug)
	}
	return removed
}

func (q *QuadTree[T]) Print(find any) {
	//fmt.Printf("qt[%v]:", q.Total)
	if q.Nodes[0] == nil {
		for i, o := range q.Objects {
			if o == find {
				fmt.Printf("QT %d. |%p| ", i, o)
			}
		}
	} else {
		q.Nodes[0].Print(find)
		q.Nodes[1].Print(find)
		q.Nodes[2].Print(find)
		q.Nodes[3].Print(find)
	}
}

func (q *QuadTree[T]) Insert(obj T) {

	q.Total++

	if q.Nodes[0] != nil {
		quadrant := q.getQuadrant(obj.bRect()) // see where it fits
		if quadrant != qDoesntFit {
			q.Nodes[quadrant].Insert(obj)
			return
		}
	}

	q.Objects = append(q.Objects, obj) // if it doesn't fit into subquadrant add it here

	if (len(q.Objects) > qtMaxObjects) && (q.Level < qtMaxLevels) {
		if q.Nodes[0] == nil {
			q.split()
		}

		i := 0
		for i < len(q.Objects) {
			quadrant := q.getQuadrant(q.Objects[i].bRect())
			if quadrant != qDoesntFit {
				objs := q.Objects[i]
				q.Objects = append(q.Objects[:i], q.Objects[i+1:]...)
				q.Nodes[quadrant].Insert(objs)
			} else {
				i++
			}
		}

	}
}

func (q *QuadTree[T]) MayCollide(r Rect) []T {
	const (
		minDist2 = PrefferredRockSize * PrefferredRockSize * 16
	)

	quadrant := q.getQuadrant(r)

	collidingObjects := q.Objects
	// for i, o := range q.Objects {
	// 	dist2 := (o.bRect().x-r.x)*(o.bRect().x-r.x) + (o.bRect().y-r.y)*(o.bRect().y-r.y)
	// 	if dist2 < minDist2 {
	// 		collidingObjects = append(collidingObjects, q.Objects[i])

	// 	}
	// }
	if q.Nodes[0] != nil {
		if quadrant != qDoesntFit {
			collidingObjects = append(collidingObjects, q.Nodes[quadrant].MayCollide(r)...)
		} else {
			for i := 0; i < 4; i++ {
				collidingObjects = append(collidingObjects, q.Nodes[i].MayCollide(r)...)
			}
		}
	}
	return collidingObjects
}
