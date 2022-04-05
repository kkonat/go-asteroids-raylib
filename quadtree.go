package main

const qtMaxObjects = 10
const qtMaxLevels = 5

type Rect struct {
	x, y int32
	w, h int32
}

type qtObj interface {
	bRect() Rect
}
type QuadTree[T qtObj] struct {
	Level   int
	Objects []T
	Bounds  Rect
	Nodes   [4]*QuadTree[T]
	Total   int
}

func NewQuadTree[T qtObj](pLevel int, pBounds Rect) *QuadTree[T] {
	return &QuadTree[T]{Level: pLevel, Bounds: pBounds}
}

func (qt *QuadTree[T]) TotalNodes() int {

	total := 0

	if qt.Nodes[0] != nil {
		for i := 0; i < len(qt.Nodes); i++ {
			total += 1
			total += qt.Nodes[i].TotalNodes()
		}
	}

	return total

}
func (q *QuadTree[T]) Clear() {

	q.Objects = nil
	for i := 0; i < 4; i++ {
		if q.Nodes[i] != nil {
			q.Nodes[i].Clear()
			q.Nodes[i] = nil
		}

	}
	q.Total = 0
}

func (q *QuadTree[T]) split() {
	subW, subH := q.Bounds.w/2, q.Bounds.h/2
	x, y := q.Bounds.x, q.Bounds.y

	q.Nodes[0] = NewQuadTree[T](q.Level+1, Rect{x, y, subW, subH})
	q.Nodes[1] = NewQuadTree[T](q.Level+1, Rect{x + subW, y, subW, subH})
	q.Nodes[2] = NewQuadTree[T](q.Level+1, Rect{x, y + subH, subW, subH})
	q.Nodes[3] = NewQuadTree[T](q.Level+1, Rect{x + subW, y + subH, subW, subH})
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
				q.Objects[i] = q.Objects[len(q.Objects)-1] // remove
				q.Objects = q.Objects[:len(q.Objects)-1]
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
	//collidingObjects := q.Objects
	collidingObjects := []T{}
	for _, o := range q.Objects {
		dist2 := (o.bRect().x-r.x)*(o.bRect().x-r.x) + (o.bRect().y-r.y)*(o.bRect().y-r.y)
		if dist2 < minDist2 {
			collidingObjects = append(collidingObjects, o)
		}
	}
	if q.Nodes[0] != nil {
		if quadrant != qDoesntFit {
			t := q.Nodes[quadrant].MayCollide(r)
			collidingObjects = append(collidingObjects, t...)
		} else {
			for i := 0; i < 4; i++ {
				collidingObjects = append(collidingObjects, q.Nodes[i].MayCollide(r)...)
			}
		}
	}
	return collidingObjects
}
