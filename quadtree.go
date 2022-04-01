package main

// https://gamedevelopment.tutsplus.com/tutorials/quick-tip-use-quadtrees-to-detect-likely-collisions-in-2d-space--gamedev-374

const maxObjects = 15
const maxLevels = 5

type Rect struct {
	x, y int32
	w, h int32
}

func (r1 Rect) equals(r2 Rect) bool {
	return r1.x == r2.x && r1.y == r2.y && r1.w == r2.w && r1.h == r2.h
}

type qtObj interface {
	bRect() Rect
}

type QuadTree struct {
	Level   int
	Objects []qtObj
	Bounds  Rect
	Nodes   [4]*QuadTree
	Total   int
}

func NewQuadTree(pLevel int, pBounds Rect) *QuadTree {
	q := new(QuadTree)
	q.Level = pLevel
	//	q.objects = make([]Obj, 50)
	q.Bounds = pBounds
	return q

}

func (qt *QuadTree) TotalNodes() int {

	total := 0

	if qt.Nodes[0] != nil {
		for i := 0; i < len(qt.Nodes); i++ {
			total += 1
			total += qt.Nodes[i].TotalNodes()
		}
	}

	return total

}
func (q *QuadTree) Clear() {

	q.Objects = nil
	for i := 0; i < 4; i++ {
		if q.Nodes[i] != nil {
			q.Nodes[i].Clear()
			q.Nodes[i] = nil
		}

	}
	q.Total = 0
}

func (q *QuadTree) split() {
	subW, subH := q.Bounds.w/2, q.Bounds.h/2
	x, y := q.Bounds.x, q.Bounds.y

	q.Nodes[0] = NewQuadTree(q.Level+1, Rect{x, y, subW, subH})
	q.Nodes[1] = NewQuadTree(q.Level+1, Rect{x + subW, y, subW, subH})
	q.Nodes[2] = NewQuadTree(q.Level+1, Rect{x, y + subH, subW, subH})
	q.Nodes[3] = NewQuadTree(q.Level+1, Rect{x + subW, y + subH, subW, subH})
}

const (
	qTopLeft     = 0
	qTopRight    = 1
	qBottomLeft  = 2
	qBottomRight = 3
	qDoesntFit   = 5
)

func (q *QuadTree) getQuadrant(o qtObj) int {

	quadrant := qDoesntFit

	cx := q.Bounds.x + (q.Bounds.w / 2)
	cy := q.Bounds.y + (q.Bounds.h / 2)

	fitsInTop := (o.bRect().y+o.bRect().h < cy)
	fitsInBottom := (o.bRect().y > cy)
	// left quadrants
	if o.bRect().x+o.bRect().w < cx {
		if fitsInTop {
			quadrant = qTopLeft
		} else if fitsInBottom {
			quadrant = qBottomLeft
		}
	} else if o.bRect().x > cx { // right quadrants
		if fitsInTop {
			quadrant = qTopRight
		} else if fitsInBottom {
			quadrant = qBottomRight
		}
	}
	return quadrant
}

func (q *QuadTree) Insert(c qtObj) {

	q.Total++

	if q.Nodes[0] != nil {
		quadrant := q.getQuadrant(c) // see where it fits
		if quadrant != qDoesntFit {
			q.Nodes[quadrant].Insert(c)
			return
		}
	}

	q.Objects = append(q.Objects, c) // if it doesn't fit into subquadrant add it here

	if (len(q.Objects) > maxObjects) && (q.Level < maxLevels) {
		if q.Nodes[0] == nil {
			q.split()
		}

		i := 0
		for i < len(q.Objects) {
			quadrant := q.getQuadrant(q.Objects[i])
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

func (q *QuadTree) MayCollide(o qtObj) []qtObj {

	var collidingObjects []qtObj

	quadrant := q.getQuadrant(o)
	if quadrant != qDoesntFit && q.Nodes[0] != nil {
		collidingObjects = append(collidingObjects, q.Nodes[quadrant].MayCollide(o)...)
	} else {
		collidingObjects = q.Objects
	}
	return collidingObjects
}
