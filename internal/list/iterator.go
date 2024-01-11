package list

type Iterator[V any] struct {
	head    *Node[V] // should be initialized in the list beginning and changes on each Next()
	current *Node[V] // current node, nil on start and changes on each after Next()
	idx     int      // current element ID
}

// Next moves iterator next.
func (iter *Iterator[V]) Next() bool {
	if iter.head == nil {
		return false
	}

	iter.current = iter.head
	iter.head = iter.head.next
	iter.idx++

	return true
}

// Idx returns the current iterator position
func (iter *Iterator[V]) Idx() int {
	return iter.idx
}

// Value returns the value
func (iter *Iterator[V]) Value() (V, bool) {
	if iter.current == nil {
		return zero[V](), false
	}

	return iter.current.value, true
}

func (iter *Iterator[V]) Node() *Node[V] {
	return iter.current
}
