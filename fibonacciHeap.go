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
		return errors.New("Negative infinity key is reserved for internal usage.")
	}

	if _, exists := heap.index[value.Tag()]; exists {
		return errors.New("Duplicate tag is not allowed")
	}

	node := new(Node)
	node.tag = value.Tag()
	node.key = value.Key()
	node.value = value
	node.children = list.New()

	heap.index[node.tag] = heap.roots.PushBack(node)
	heap.num++

	if heap.min == nil || heap.min.Value.(*Node).key > node.key {
		heap.min = heap.index[node.tag]
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

	min := heap.min

	children := heap.min.Value.(*Node).children
	if children != nil {
		for e := children.Front(); e != nil; e = e.Next() {
			node := e.Value.(*Node)
			node.parent = nil
			heap.index[node.tag] = heap.roots.PushBack(node)
			for child := node.children.Front(); child != nil; child = child.Next() {
				child.Value.(*Node).parent = heap.index[node.tag]
			}
		}
	}

	heap.roots.Remove(heap.min)
	delete(heap.index, heap.min.Value.(*Node).tag)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.min = min.Next()
	}

	heap.consolidate()

	return min.Value.(*Node).value
}

func (heap *fibHeap) Union(another FibHeap) error {
	anotherHeap, safe := another.(*fibHeap)
	if !safe {
		return errors.New("Target heap is not a valid Fibonacci Heap")
	}
	for tag, _ := range anotherHeap.index {
		if _, exists := heap.index[tag]; exists {
			return errors.New("Duplicate tag is found in the target heap")
		}
	}

	for _, element := range anotherHeap.index {
		heap.Insert(element.Value.(*Node).value)
	}

	return nil
}

func (heap *fibHeap) DecreaseKey(value Value) error {
	var element *list.Element
	var exists bool
	if element, exists = heap.index[value.Tag()]; !exists {
		return errors.New("Value is not found")
	}

	return heap.decreaseKey(element, value, value.Key())
}

func (heap *fibHeap) decreaseKey(element *list.Element, value Value, key float64) error {
	if key > element.Value.(*Node).key {
		return errors.New("New key is greater than current key")
	}

	node := element.Value.(*Node)
	node.key = key
	node.value = value
	if node.parent != nil {
		parent := node.parent
		if node.key < node.parent.Value.(*Node).key {
			element = heap.cut(element)
			heap.cascadingCut(parent)
		}
	}

	if node.parent == nil && node.key < heap.min.Value.(*Node).key {
		heap.min = element
	}

	return nil
}

func (heap *fibHeap) Delete(value Value) error {
	if _, exists := heap.index[value.Tag()]; !exists {
		return errors.New("Value is not found")
	}

	heap.ExtractTag(value.Tag())

	return nil
}

func (heap *fibHeap) GetTag(tag interface{}) (value Value) {
	if element, exists := heap.index[tag]; exists {
		value = element.Value.(*Node).value
	}

	return
}

func (heap *fibHeap) ExtractTag(tag interface{}) (value Value) {
	var element *list.Element
	var exists bool
	if element, exists = heap.index[tag]; !exists {
		return
	}

	value = element.Value.(*Node).value
	heap.decreaseKey(element, element.Value.(*Node).value, math.Inf(-1))
	heap.ExtractMin()

	return
}

func (heap *fibHeap) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Total number: %d, Root Size: %d, Index size: %d,\n", heap.num, heap.roots.Len(), len(heap.index)))
	buffer.WriteString(fmt.Sprintf("Current minimun: key(%f), tag(%v), value(%v),\n", heap.min.Value.(*Node).key, heap.min.Value.(*Node).tag, heap.min.Value.(*Node).value))
	buffer.WriteString(fmt.Sprintf("Heap detail: "))
	probeTree(&buffer, heap.roots)

	return buffer.String()
}

func probeTree(buffer *bytes.Buffer, tree *list.List) {
	buffer.WriteString(fmt.Sprintf("< "))
	for e := tree.Front(); e != nil; e = e.Next() {
		buffer.WriteString(fmt.Sprintf("%f ", e.Value.(*Node).key))
		children := e.Value.(*Node).children
		if children.Len() != 0 {
			probeTree(buffer, children)
		}
	}
	buffer.WriteString(fmt.Sprintf("> "))
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
	heap.index[child.Value.(*Node).tag] = parent.Value.(*Node).children.PushBack(child.Value)
	for grandChild := child.Value.(*Node).children.Front(); grandChild != nil; grandChild = grandChild.Next() {
		grandChild.Value.(*Node).parent = heap.index[child.Value.(*Node).tag]
	}
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

func (heap *fibHeap) cut(element *list.Element) *list.Element {
	node := element.Value.(*Node)
	node.parent.Value.(*Node).children.Remove(element)
	node.parent.Value.(*Node).degree--
	node.parent = nil
	node.marked = false
	heap.index[node.tag] = heap.roots.PushBack(node)
	for child := node.children.Front(); child != nil; child = child.Next() {
		child.Value.(*Node).parent = heap.index[node.tag]
	}

	return heap.index[node.tag]
}

func (heap *fibHeap) cascadingCut(element *list.Element) {
	node := element.Value.(*Node)
	if node.parent != nil {
		if !node.marked {
			node.marked = true
		} else {
			parent := node.parent
			heap.cut(element)
			heap.cascadingCut(parent)
		}
	}
}
