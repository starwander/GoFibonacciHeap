// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

// Package fibHeap implements the Fibonacci Heap priority queue.
// This implementation is a bit different from the traditional Fibonacci Heap by having an index map to achieve better encapsulation.
package fibHeap

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math"
)

// Value is the interface that all values push into or pop from the FibHeap must implement.
type Value interface {
	// Tag returns the unique tag of the value.
	// The tag is used in the index map.
	Tag() interface{}
	// Key returns the key as known as the priority of the value.
	// The valid range of the key is (-inf, +inf].
	Key() float64
}

// FibHeap represents a Fibonacci Heap.
// Please note that all methods of FibHeap are not concurrent safe.
type FibHeap struct {
	roots       *list.List
	index       map[interface{}]*node
	treeDegrees map[uint]*list.Element
	min         *node
	num         uint
}

type node struct {
	self     *list.Element
	parent   *node
	children *list.List
	marked   bool
	degree   uint
	position uint
	tag      interface{}
	key      float64
	value    Value
}

// NewFibHeap creates an initialized Fibonacci Heap.
func NewFibHeap() *FibHeap {
	heap := new(FibHeap)
	heap.roots = list.New()
	heap.index = make(map[interface{}]*node)
	heap.treeDegrees = make(map[uint]*list.Element)
	heap.num = 0
	heap.min = nil

	return heap
}

// Num returns the total number of values in the heap.
func (heap *FibHeap) Num() uint {
	return heap.num
}

// Insert pushes the input value into the heap.
// The input value must implements the Value interface.
// Try to insert a duplicate tag value will cause an error return.
// The valid range of the value's key is (-inf, +inf].
// Try to insert a -inf key value will cause an error return.
// Insert will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) Insert(value Value) error {
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

// Minimum returns the current minimum value in the heap by key.
// Minimum will not extract the value so the value will still exists in the heap.
// An empty heap will return nil.
func (heap *FibHeap) Minimum() Value {
	if heap.num == 0 {
		return nil
	}

	return heap.min.value
}

// ExtractMin returns the current minimum value in the heap and then extracts the value from the heap.
// An empty heap will return nil and extracts nothing.
func (heap *FibHeap) ExtractMin() Value {
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
	heap.treeDegrees[min.position] = nil
	delete(heap.index, heap.min.tag)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.consolidate()
	}

	return min.value
}

// Union merges the input heap in.
// All values of the input heap must not have duplicate tags. Otherwise an error will be returned.
func (heap *FibHeap) Union(anotherHeap *FibHeap) error {
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

// DecreaseKey updates the value in the heap by the input value.
// If the input value has a larger key or -inf key, an error will be returned.
// If the tag of the input value is not existed in the heap, an error will be returned.
// DecreaseKey will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) DecreaseKey(value Value) error {
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

// Delete deletes the value in the heap by the input value.
// If the tag of the input value is not existed in the heap, an error will be returned.
// Delete will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) Delete(value Value) error {
	if value == nil {
		return errors.New("Input value is nil.")
	}

	if _, exists := heap.index[value.Tag()]; !exists {
		return errors.New("Value is not found")
	}

	heap.ExtractTag(value.Tag())

	return nil
}

// GetTag searches and returns the value in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// GetTag will not extract the value so the value will still exist in the heap.
func (heap *FibHeap) GetTag(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
	}

	return
}

// ExtractTag searches and extracts the value in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// ExtractTag will extract the value so the value will no longer exist in the heap.
func (heap *FibHeap) ExtractTag(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		heap.decreaseKey(node, node.value, math.Inf(-1))
		heap.ExtractMin()
		value = node.value
		return
	}

	return
}

// String provides some basic debug information of the heap.
// It returns the total number, roots size, index size and current minimum value of the heap.
// It also returns the topology of the trees by dfs search.
func (heap *FibHeap) String() string {
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

func (heap *FibHeap) consolidate() {
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		heap.treeDegrees[tree.Value.(*node).position] = nil
	}

	for tree := heap.roots.Front(); tree != nil; {
		if heap.treeDegrees[tree.Value.(*node).degree] == nil {
			heap.treeDegrees[tree.Value.(*node).degree] = tree
			tree.Value.(*node).position = tree.Value.(*node).degree
			tree = tree.Next()
			continue
		}

		if heap.treeDegrees[tree.Value.(*node).degree] == tree {
			tree = tree.Next()
			continue
		}

		for heap.treeDegrees[tree.Value.(*node).degree] != nil {
			anotherTree := heap.treeDegrees[tree.Value.(*node).degree]
			heap.treeDegrees[tree.Value.(*node).degree] = nil
			if tree.Value.(*node).key <= anotherTree.Value.(*node).key {
				heap.roots.Remove(anotherTree)
				heap.link(tree.Value.(*node), anotherTree.Value.(*node))
			} else {
				heap.roots.Remove(tree)
				heap.link(anotherTree.Value.(*node), tree.Value.(*node))
				tree = anotherTree
			}
		}
		heap.treeDegrees[tree.Value.(*node).degree] = tree
		tree.Value.(*node).position = tree.Value.(*node).degree
	}

	heap.resetMin()
}

func (heap *FibHeap) link(parent, child *node) {
	child.marked = false
	child.parent = parent
	child.self = parent.children.PushBack(child)
	parent.degree++
}

func (heap *FibHeap) resetMin() {
	key := math.Inf(1)
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		if tree.Value.(*node).key < key {
			heap.min = tree.Value.(*node)
			key = tree.Value.(*node).key
		}
	}
}

func (heap *FibHeap) decreaseKey(node *node, value Value, key float64) error {
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

func (heap *FibHeap) cut(node *node) {
	node.parent.children.Remove(node.self)
	node.parent.degree--
	node.parent = nil
	node.marked = false
	node.self = heap.roots.PushBack(node)
}

func (heap *FibHeap) cascadingCut(node *node) {
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
