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
