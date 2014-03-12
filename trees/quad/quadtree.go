package quad

type Positionable interface {
	X() float32
	Y() float32
}

type quadNode struct {
	// the smaller subdivisions of this block
	// len(children) == 0 if does not contain
	// smaller divisions
	children    []*quadNode
	numChildren int

	// the data points stored in this quadrant
	data [][]Positionable

	// the position of this node in the overall world
	// (offset from root (0, 0))
	bottomLeftX float32
	bottomLeftY float32
	topRightX   float32
	topRightY   float32

	// might make more sense to keep track of four
	// different slices of data, so when we need to divide
	// we don't have to look over entire data slice
}

type QuadTree struct {
	// The root of the quadtree
	root *quadNode

	// Limit of number of items a division can contain before it needs to be subdivided
	divisionLimit int
}

var (
	divisionLimit    int
	maxXVal, maxYVal float32
)

func NewQuadTree(divLimit int, maxX, maxY float32) *QuadTree {
	divisionLimit = divLimit
	return &QuadTree{
		root: &quadNode{
			children: []*quadNode{
				&quadNode{ // top right
					data:        make([][]Positionable, 4),
					bottomLeftY: 0,
					bottomLeftX: 0,
					topRightY:   maxY,
					topRightX:   maxX,
				},
				&quadNode{ // top left
					data:        make([][]Positionable, 4),
					bottomLeftY: 0,
					bottomLeftX: maxX * -1,
					topRightY:   maxY,
					topRightX:   0,
				},
				&quadNode{ // bottom left
					data:        make([][]Positionable, 4),
					bottomLeftY: maxY * -1,
					bottomLeftX: maxX * -1,
					topRightY:   0,
					topRightX:   0,
				},
				&quadNode{ // bottom right
					data:        make([][]Positionable, 4),
					bottomLeftY: maxY * -1,
					bottomLeftX: 0,
					topRightY:   0,
					topRightX:   maxX,
				},
			},
			bottomLeftY: maxY * -1,
			bottomLeftX: maxX * -1,
			topRightY:   maxY,
			topRightX:   maxX,
		},
	}
}

func (q *QuadTree) Insert(p Positionable) {
	q.root.insert(p)
}

func (q *QuadTree) Delete(p Positionable) error {
	// TODO(ttacon)
	return nil
}

func (q *QuadTree) Within(xRange, yRange float32, p Positionable, suitcase interface{}) error {
	// TODO(ttacon)
	return nil
}

func (q *quadNode) insert(p Positionable) {
	if len(q.children) > 0 {
		// find correct child and return its insert
		if q.children[0].contains(p) {
			q.children[0].insert(p)
			return
		}
		if q.children[1].contains(p) {
			q.children[1].insert(p)
			return
		}
		if q.children[2].contains(p) {
			q.children[2].insert(p)
			return
		}
		if q.children[3].contains(p) {
			q.children[3].insert(p)
			return
		}
	}

	// if we get here we can just add the node
	// TODO(ttacon): these should honestly be precalculated
	var (
		midX, midY float32
	)
	if q.bottomLeftX > 0 {
		midX = ((q.topRightX - q.bottomLeftX) / 2) + q.bottomLeftX
	} else {
		midX = q.topRightX + ((q.bottomLeftX - q.topRightX) / 2)
	}
	if q.bottomLeftY > 0 {
		midY = ((q.topRightY - q.bottomLeftY) / 2) + q.bottomLeftY
	} else {
		midY = q.topRightY + ((q.bottomLeftY - q.topRightY) / 2)
	}

	if midX < p.X() {
		if midY < p.Y() {
			q.data[0] = append(q.data[0], p)
		} else {
			q.data[3] = append(q.data[3], p)
		}
	} else {
		if midY < p.Y() {
			q.data[1] = append(q.data[1], p)
		} else {
			q.data[2] = append(q.data[2], p)
		}
	}
	q.numChildren += 1

	// check node isn't full
	if q.numChildren == divisionLimit {
		q.split()
	}
}

func (q *quadNode) contains(p Positionable) bool {
	return q.bottomLeftY <= p.Y() && q.topRightY >= p.Y() &&
		q.bottomLeftX <= p.X() && q.topRightX >= p.X()
}

func (q *quadNode) split() {
	var (
		midX, midY float32
	)
	if q.bottomLeftX > 0 {
		midX = ((q.topRightX - q.bottomLeftX) / 2) + q.bottomLeftX
	} else {
		midX = q.topRightX + ((q.bottomLeftX - q.topRightX) / 2)
	}
	if q.bottomLeftY > 0 {
		midY = ((q.topRightY - q.bottomLeftY) / 2) + q.bottomLeftY
	} else {
		midY = q.topRightY + ((q.bottomLeftY - q.topRightY) / 2)
	}

	q.children = []*quadNode{
		&quadNode{ // quadrant 0
			data:        make([][]Positionable, 4),
			bottomLeftX: midX,
			bottomLeftY: midY,
			topRightX:   q.topRightX,
			topRightY:   q.topRightY,
		},
		&quadNode{ // quadrant 1
			data:        make([][]Positionable, 4),
			bottomLeftX: q.bottomLeftX,
			bottomLeftY: midY,
			topRightX:   midX,
			topRightY:   q.topRightY,
		},
		&quadNode{ // quadrant 2
			data:        make([][]Positionable, 4),
			bottomLeftX: q.bottomLeftX,
			bottomLeftY: q.bottomLeftY,
			topRightX:   midX,
			topRightY:   midY,
		},
		&quadNode{ // quadrant 3
			data:        make([][]Positionable, 4),
			bottomLeftX: midX,
			bottomLeftY: q.bottomLeftY,
			topRightX:   q.topRightX,
			topRightY:   midY,
		},
	}
	q.data = nil
	q.numChildren = 0
}
