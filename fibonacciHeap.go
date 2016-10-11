package fibHeap

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math"
)

type fibHeap struct {
	roots *list.List
	index map[interface{}]*node
	min   *node
	num   uint
}

type node struct {
	self     *list.Element
	parent   *node
	children *list.List
	marked   bool
	degree   uint
	tag      interface{}
	key      float64
	value    Value
}

func (heap *fibHeap) Num() uint {
	return heap.num
}

func (heap *fibHeap) Insert(value Value) error {
	if value == nil {
		return errors.New("Input value is nil.")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage.")
	}

	if _, exists := heap.index[value.Tag()]; exists {
		return errors.New("Duplicate tag is not allowed")
	}

	node := new(node)
	node.children = list.New()
	node.tag = value.Tag()
	node.key = value.Key()
	node.value = value

	node.self = heap.roots.PushBack(node)
	heap.index[node.tag] = node
	heap.num++

	if heap.min == nil || heap.min.key > node.key {
		heap.min = node
	}

	return nil
}

func (heap *fibHeap) Minimum() Value {
	if heap.num == 0 {
		return nil
	}

	return heap.min.value
}

func (heap *fibHeap) ExtractMin() Value {
	if heap.num == 0 {
		return nil
	}

	min := heap.min

	children := heap.min.children
	if children != nil {
		for e := children.Front(); e != nil; e = e.Next() {
			e.Value.(*node).parent = nil
			e.Value.(*node).self = heap.roots.PushBack(e.Value.(*node))
		}
	}

	heap.roots.Remove(heap.min.self)
	delete(heap.index, heap.min.tag)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.consolidate()
	}

	return min.value
}

func (heap *fibHeap) Union(another FibHeap) error {
	anotherHeap, safe := another.(*fibHeap)
	if !safe {
		return errors.New("Target heap is not a valid Fibonacci Heap")
	}

	for tag := range anotherHeap.index {
		if _, exists := heap.index[tag]; exists {
			return errors.New("Duplicate tag is found in the target heap")
		}
	}

	for _, node := range anotherHeap.index {
		heap.Insert(node.value)
	}

	return nil
}

func (heap *fibHeap) DecreaseKey(value Value) error {
	if value == nil {
		return errors.New("Input value is nil.")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage.")
	}

	if node, exists := heap.index[value.Tag()]; exists {
		return heap.decreaseKey(node, value, value.Key())
	}

	return errors.New("Value is not found")
}

func (heap *fibHeap) Delete(value Value) error {
	if value == nil {
		return errors.New("Input value is nil.")
	}

	if _, exists := heap.index[value.Tag()]; !exists {
		return errors.New("Value is not found")
	}

	heap.ExtractTag(value.Tag())

	return nil
}

func (heap *fibHeap) GetTag(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
	}

	return
}

func (heap *fibHeap) ExtractTag(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		heap.decreaseKey(node, node.value, math.Inf(-1))
		heap.ExtractMin()
		value = node.value
		return
	}

	return
}

func (heap *fibHeap) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Total number: %d, Root Size: %d, Index size: %d,\n", heap.num, heap.roots.Len(), len(heap.index)))
	buffer.WriteString(fmt.Sprintf("Current minimun: key(%f), tag(%v), value(%v),\n", heap.min.key, heap.min.tag, heap.min.value))
	buffer.WriteString(fmt.Sprintf("Heap detail:\n"))
	probeTree(&buffer, heap.roots)
	buffer.WriteString(fmt.Sprintf("\n"))

	return buffer.String()
}

func probeTree(buffer *bytes.Buffer, tree *list.List) {
	buffer.WriteString(fmt.Sprintf("< "))
	for e := tree.Front(); e != nil; e = e.Next() {
		buffer.WriteString(fmt.Sprintf("%f ", e.Value.(*node).key))
		if e.Value.(*node).children.Len() != 0 {
			probeTree(buffer, e.Value.(*node).children)
		}
	}
	buffer.WriteString(fmt.Sprintf("> "))
}

func (heap *fibHeap) consolidate() {
	treeDegrees := make(map[uint]*list.Element, heap.maxPossibleNum())

	for tree := heap.roots.Front(); tree != nil; {
		degree := tree.Value.(*node).degree

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
			if tree.Value.(*node).key <= anotherTree.Value.(*node).key {
				heap.roots.Remove(anotherTree)
				heap.link(tree.Value.(*node), anotherTree.Value.(*node))
			} else {
				heap.roots.Remove(tree)
				heap.link(anotherTree.Value.(*node), tree.Value.(*node))
				tree = anotherTree
			}
			degree++
		}
		treeDegrees[degree] = tree
	}

	heap.resetMin()
}

func (heap *fibHeap) maxPossibleNum() int {
	maxPossibleRootNum := (int)(math.Ceil(-1 + math.Sqrt(float64(1+8*heap.num))/2))
	if heap.roots.Len() < maxPossibleRootNum {
		return heap.roots.Len()
	}

	return maxPossibleRootNum
}

func (heap *fibHeap) link(parent, child *node) {
	child.marked = false
	child.parent = parent
	child.self = parent.children.PushBack(child)
	parent.degree++
}

func (heap *fibHeap) resetMin() {
	key := math.Inf(1)
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		if tree.Value.(*node).key < key {
			heap.min = tree.Value.(*node)
			key = tree.Value.(*node).key
		}
	}
}

func (heap *fibHeap) decreaseKey(node *node, value Value, key float64) error {
	if key > node.key {
		return errors.New("New key is greater than current key")
	}

	node.key = key
	node.value = value
	if node.parent != nil {
		parent := node.parent
		if node.key < node.parent.key {
			heap.cut(node)
			heap.cascadingCut(parent)
		}
	}

	if node.parent == nil && node.key < heap.min.key {
		heap.min = node
	}

	return nil
}

func (heap *fibHeap) cut(node *node) {
	node.parent.children.Remove(node.self)
	node.parent.degree--
	node.parent = nil
	node.marked = false
	node.self = heap.roots.PushBack(node)
}

func (heap *fibHeap) cascadingCut(node *node) {
	if node.parent != nil {
		if !node.marked {
			node.marked = true
		} else {
			parent := node.parent
			heap.cut(node)
			heap.cascadingCut(parent)
		}
	}
}
