package list

type Node[V any] struct {
	next  *Node[V]
	value V
}

func (n *Node[V]) Value() V {
	return n.value
}

func (n *Node[V]) Next() *Node[V] {
	return n.next
}

// List implements a singly linked list.
type List[V any] struct {
	len  int
	head *Node[V]
}

func New[V any](elements ...V) List[V] {
	l := List[V]{}
	for _, v := range elements {
		l.Push(v)
	}

	return l
}

func (l List[V]) Iter() Iterator[V] {
	return Iterator[V]{head: l.head, idx: -1}
}

// CollectValues yields the list as a slice of values.
// The list cannot be mutated through the slice.
func (l *List[V]) CollectValues() []V {
	if l.len == 0 {
		return nil
	}

	var (
		slice = make([]V, 0, l.len)
		iter  = l.Iter()
	)

	for iter.Next() {
		v, _ := iter.Value()
		slice = append(slice, v)
	}

	return slice
}

// Unshift inserts value in front of the list with O(1) complexity.
func (l *List[V]) Unshift(v V) {
	l.head = &Node[V]{value: v, next: l.head}
	l.len++
}

// Push value to the end with O(n) complexity.
func (l *List[V]) Push(value V) {
	defer func() { l.len++ }()

	var (
		iter = l.Iter()
		node = &Node[V]{value: value}
	)

	if l.head == nil {
		l.head = node
		return
	}

	// seek to the end
	for iter.Next() {
	}

	iter.Node().next = node
}

// Nth returns an N-th node starting from the head.
func (l *List[V]) Nth(n int) *Node[V] {
	if n >= l.len {
		return nil
	}

	var iter = l.Iter()
	for iter.Next() {
		if iter.Idx() == n {
			return iter.Node()
		}
	}

	return nil
}

// Delete n-th element and return true if the element existed.
func (l *List[V]) Delete(n int) bool {
	if l.head == nil {
		return false
	}

	if l.head != nil && n == 0 {
		l.head = l.Nth(1)
		return true
	}

	prev := l.Nth(n - 1)
	if prev == nil {
		return false
	}

	prev.next = l.Nth(n + 1)

	return true
}
