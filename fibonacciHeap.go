package fibHeap

import (
	"container/list"
	"errors"
	"math"
)

type fibHeap struct {
	roots *list.List
	index map[interface{}]*list.Element
	min   *list.Element
	num   uint
}

type Node struct {
	children *list.List
	parent   *list.Element
	tag      interface{}
	key      float64
	value    Value
	marked   bool
	degree   uint
}

func (heap *fibHeap) Num() uint {
	return heap.num
}

func (heap *fibHeap) Insert(value Value) error {
	if value.Key() <= math.Inf(-1) {
		return errors.New("Negative infinity is reserved for internal usage.")
	}

	if _, exists := heap.index[value.Tag()]; exists {
		return errors.New("Duplicate tag is not allowed")
	}

	node := new(Node)
	node.tag = value.Tag()
	node.key = value.Key()
	node.value = value
	node.children = list.New()

	element := heap.roots.PushBack(node)
	heap.index[node.tag] = element
	heap.num++

	if heap.min == nil || heap.min.Value.(*Node).key > node.key {
		heap.min = element
	}

	return nil
}

func (heap *fibHeap) Minimum() Value {
	if heap.num == 0 {
		return nil
	}

	return heap.min.Value.(*Node).value
}

func (heap *fibHeap) ExtractMin() Value {
	if heap.num == 0 {
		return nil
	}

	min := heap.min.Value

	children := heap.min.Value.(*Node).children
	if children != nil {
		for e := children.Front(); e != nil; e = e.Next() {
			e.Value.(*Node).parent = nil

		}
		heap.roots.PushBackList(children)
	}

	heap.roots.Remove(heap.min)
	delete(heap.index, heap.min.Value.(*Node).tag)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.min = heap.min.Next()
	}

	heap.consolidate()

	return min.(*Node).value
}

func (heap *fibHeap) Union(another FibHeap) error {
	anotherHeap, safe := another.(*fibHeap)
	if !safe {
		return nil
	}
	for tag, _ := range anotherHeap.index {
		if _, exists := heap.index[tag]; exists {
			return errors.New("Duplicate tag is found in the target heap")
		}
	}

	heap.roots.PushBackList(anotherHeap.roots)
	for tag, element := range anotherHeap.index {
		heap.index[tag] = element
	}
	heap.num += anotherHeap.num
	if heap.min == nil || (anotherHeap.min != nil && anotherHeap.min.Value.(*Node).key < heap.min.Value.(*Node).key) {
		heap.min = anotherHeap.min
	}

	return nil
}

func (heap *fibHeap) GetTag(tag interface{}) (value Value) {
	if element, exists := heap.index[tag]; exists {
		value = element.Value.(*Node).value
	}

	return
}

func (heap *fibHeap) consolidate() {
	if heap.num == 0 {
		return
	}

	treeDegrees := make(map[uint]*list.Element)

	for tree := heap.roots.Front(); tree != nil; {
		degree := tree.Value.(*Node).degree

		if treeDegrees[degree] == nil {
			treeDegrees[degree] = tree
			tree = tree.Next()
			continue
		}

		if treeDegrees[degree] == tree {
			tree = tree.Next()
			continue
		}

		for treeDegrees[degree] != nil {
			anotherTree := treeDegrees[degree]
			treeDegrees[degree] = nil
			if tree.Value.(*Node).key <= anotherTree.Value.(*Node).key {
				heap.roots.Remove(anotherTree)
				heap.link(tree, anotherTree)
			} else {
				heap.roots.Remove(tree)
				heap.link(anotherTree, tree)
				tree = anotherTree
			}
			degree++
		}
		treeDegrees[degree] = tree
	}

	heap.resetMin()
}

func (heap *fibHeap) link(parent, child *list.Element) {
	child.Value.(*Node).marked = false
	child.Value.(*Node).parent = parent
	parent.Value.(*Node).children.PushBack(child.Value)
	parent.Value.(*Node).degree++
}

func (heap *fibHeap) resetMin() {
	key := math.Inf(1)
	heap.min = nil
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		if tree.Value.(*Node).key < key {
			heap.min = tree
			key = tree.Value.(*Node).key
		}
	}
}
