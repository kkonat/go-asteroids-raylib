package quadtree

import (
	"fmt"
)

const qtMaxObjects = 10
const qtMaxLevels = 5

type Rect struct {
	X, Y int32
	W, H int32
}

type QtObj interface {
	comparable
	BRect() Rect
}

type QuadTree[T QtObj] struct {
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
func NewNode[T QtObj](pLevel int, pBounds Rect) *QuadTree[T] {
	n := new(QuadTree[T])
	n.Objects = make([]T, 0, qtMaxObjects)
	n.Level = pLevel
	n.Bounds = pBounds
	return n
}
func (q *QuadTree[T]) Split() {
	subW, subH := q.Bounds.W/2, q.Bounds.H/2
	x, y := q.Bounds.X, q.Bounds.Y

	q.Nodes[0] = NewNode[T](q.Level+1, Rect{x, y, subW, subH})
	q.Nodes[1] = NewNode[T](q.Level+1, Rect{x + subW, y, subW, subH})
	q.Nodes[2] = NewNode[T](q.Level+1, Rect{x, y + subH, subW, subH})
	q.Nodes[3] = NewNode[T](q.Level+1, Rect{x + subW, y + subH, subW, subH})
}

const (
	qTopLeft     = 0
	qTopRight    = 1
	qBottomLeft  = 2
	qBottomRight = 3
	qDoesntFit   = 5
)

// checks quadrant of a tree element
func (q *QuadTree[T]) getQuadrant(r Rect) int {

	quadrant := qDoesntFit

	cx := q.Bounds.X + (q.Bounds.W / 2)
	cy := q.Bounds.Y + (q.Bounds.H / 2)

	fitsInTop := (r.Y+r.H < cy)
	fitsInBottom := (r.Y > cy)
	// left quadrants
	if r.X+r.W < cx {
		if fitsInTop {
			quadrant = qTopLeft
		} else if fitsInBottom {
			quadrant = qBottomLeft
		}
	} else if r.X > cx { // right quadrants
		if fitsInTop {
			quadrant = qTopRight
		} else if fitsInBottom {
			quadrant = qBottomRight
		}
	}
	return quadrant
}

// Removes an element from the quadree
func (q *QuadTree[T]) Remove(objTbRemved T) bool {
	removed := false
	if q.Nodes[0] == nil {
		for i, o := range q.Objects {
			if o == objTbRemved {
				var zeroValue T
				q.Objects[i] = zeroValue // hack to assign a zero value of an arbitrary type T to a tree node
				q.Objects = append(q.Objects[:i], q.Objects[i+1:]...)
				q.Total--
				removed = true
				break
			}
		}
	} else {
		removed = removed || q.Nodes[0].Remove(objTbRemved)
		removed = removed || q.Nodes[1].Remove(objTbRemved)
		removed = removed || q.Nodes[2].Remove(objTbRemved)
		removed = removed || q.Nodes[3].Remove(objTbRemved)
	}
	if !removed {
		fmt.Print("x")
	}
	return removed
}

// inserts an element into the quadree
func (q *QuadTree[T]) Insert(obj T) {

	q.Total++

	if q.Nodes[0] != nil {
		quadrant := q.getQuadrant(obj.BRect()) // see where it fits
		if quadrant != qDoesntFit {
			q.Nodes[quadrant].Insert(obj)
			return
		}
	}

	q.Objects = append(q.Objects, obj) // if it doesn't fit into subquadrant add it here

	if (len(q.Objects) > qtMaxObjects) && (q.Level < qtMaxLevels) {
		if q.Nodes[0] == nil {
			q.Split()
		}

		i := 0
		for i < len(q.Objects) {
			quadrant := q.getQuadrant(q.Objects[i].BRect())
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

// returns a slice of objects that may collide with the given Rectangle AND which are within the minimum distance squared
func (q *QuadTree[T]) MayCollide(r Rect, miniDist2 int32) []T {

	quadrant := q.getQuadrant(r)

	collidingObjects := make([]T, 0) // empty slice

	for i := range q.Objects {
		collidingObjects = append(collidingObjects, q.Objects[i])
	}
	if q.Nodes[0] != nil {
		if quadrant != qDoesntFit {
			collidingObjects = append(collidingObjects, q.Nodes[quadrant].MayCollide(r, miniDist2)...)
		} else {

			for i := 0; i < 4; i++ {
				collidingObjects = append(collidingObjects, q.Nodes[i].MayCollide(r, miniDist2)...)

			}
		}
	}
	return collidingObjects
}
