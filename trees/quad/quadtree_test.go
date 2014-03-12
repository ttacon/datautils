package quad

import "testing"

type P struct {
	val int
	x   float32
	y   float32
}

func (p P) X() float32 { return p.x }
func (p P) Y() float32 { return p.y }

func BenchmarkInsert(b *testing.B) {
	q := NewQuadTree(20, 100, 100)
	for i := 0; i < b.N; i++ {
		q.Insert(P{i, float32(i % 100), float32(i % 100)})
	}
}
