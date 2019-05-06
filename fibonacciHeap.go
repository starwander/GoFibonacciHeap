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

// Value is the interface that all values push into or pop from the FibHeap by value interfaces must implement.
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

// Insert pushes the input tag and key into the heap.
// Try to insert a duplicate tag value will cause an error return.
// The valid range of the key is (-inf, +inf].
// Try to insert a -inf key value will cause an error return.
// Insert will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) Insert(tag interface{}, key float64) error {
	if tag == nil {
		return errors.New("Input tag is nil ")
	}

	return heap.insert(tag, key, nil)
}

// InsertValue pushes the input value into the heap.
// The input value must implements the Value interface.
// Try to insert a duplicate tag value will cause an error return.
// The valid range of the value's key is (-inf, +inf].
// Try to insert a -inf key value will cause an error return.
// Insert will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) InsertValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	return heap.insert(value.Tag(), value.Key(), value)
}

// Minimum returns the current minimum tag and key in the heap sorted by the key.
// Minimum will not extract the tag and key so the value will still exists in the heap.
// An empty heap will return nil and -inf.
func (heap *FibHeap) Minimum() (interface{}, float64) {
	if heap.num == 0 {
		return nil, math.Inf(-1)
	}

	return heap.min.tag, heap.min.key
}

// MinimumValue returns the current minimum value in the heap sorted by the key.
// MinimumValue will not extract the value so the value will still exists in the heap.
// An empty heap will return nil.
func (heap *FibHeap) MinimumValue() Value {
	if heap.num == 0 {
		return nil
	}

	return heap.min.value
}

// ExtractMin returns the current minimum tag and key in the heap and then extracts them from the heap.
// An empty heap will return nil/-inf and extracts nothing.
func (heap *FibHeap) ExtractMin() (interface{}, float64) {
	if heap.num == 0 {
		return nil, math.Inf(-1)
	}

	min := heap.extractMin()

	return min.tag, min.key
}

// ExtractMinValue returns the current minimum value in the heap and then extracts it from the heap.
// An empty heap will return nil and extracts nothing.
func (heap *FibHeap) ExtractMinValue() Value {
	if heap.num == 0 {
		return nil
	}

	min := heap.extractMin()

	return min.value
}

// Union merges the input heap in.
// All values of the input heap must not have duplicate tags. Otherwise an error will be returned.
func (heap *FibHeap) Union(anotherHeap *FibHeap) error {
	for tag := range anotherHeap.index {
		if _, exists := heap.index[tag]; exists {
			return errors.New("Duplicate tag is found in the target heap ")
		}
	}

	for _, node := range anotherHeap.index {
		heap.InsertValue(node.value)
	}

	return nil
}

// DecreaseKey updates the tag in the heap by the input key.
// If the input key has a larger key or -inf key, an error will be returned.
// If the input tag is not existed in the heap, an error will be returned.
// DecreaseKey will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) DecreaseKey(tag interface{}, key float64) error {
	if tag == nil {
		return errors.New("Input tag is nil ")
	}

	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[tag]; exists {
		return heap.decreaseKey(node, nil, key)
	}

	return errors.New("Value is not found ")
}

// DecreaseKeyValue updates the value in the heap by the input value.
// If the input value has a larger key or -inf key, an error will be returned.
// If the tag of the input value is not existed in the heap, an error will be returned.
// DecreaseKeyValue will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) DecreaseKeyValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[value.Tag()]; exists {
		return heap.decreaseKey(node, value, value.Key())
	}

	return errors.New("Value is not found ")
}

// IncreaseKey updates the tag in the heap by the input key.
// If the input key has a smaller key or -inf key, an error will be returned.
// If the input tag is not existed in the heap, an error will be returned.
// IncreaseKey will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) IncreaseKey(tag interface{}, key float64) error {
	if tag == nil {
		return errors.New("Input tag is nil ")
	}

	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[tag]; exists {
		return heap.increaseKey(node, nil, key)
	}

	return errors.New("Value is not found ")
}

// IncreaseKeyValue updates the value in the heap by the input value.
// If the input value has a smaller key or -inf key, an error will be returned.
// If the tag of the input value is not existed in the heap, an error will be returned.
// IncreaseKeyValue will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) IncreaseKeyValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	if math.IsInf(value.Key(), -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if node, exists := heap.index[value.Tag()]; exists {
		return heap.increaseKey(node, value, value.Key())
	}

	return errors.New("Value is not found ")
}

// Delete deletes the input tag in the heap.
// If the input tag is not existed in the heap, an error will be returned.
// Delete will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) Delete(tag interface{}) error {
	if tag == nil {
		return errors.New("Input tag is nil ")
	}

	if _, exists := heap.index[tag]; !exists {
		return errors.New("Tag is not found ")
	}

	heap.ExtractValue(tag)

	return nil
}

// DeleteValue deletes the value in the heap by the input value.
// If the tag of the input value is not existed in the heap, an error will be returned.
// DeleteValue will check the nil interface but not the interface with nil value.
// Try to input of an interface with nil value will cause invalid address panic.
func (heap *FibHeap) DeleteValue(value Value) error {
	if value == nil {
		return errors.New("Input value is nil ")
	}

	if _, exists := heap.index[value.Tag()]; !exists {
		return errors.New("Value is not found ")
	}

	heap.ExtractValue(value.Tag())

	return nil
}

// GetTag searches and returns the key in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// GetTag will not extract the value so the value will still exist in the heap.
func (heap *FibHeap) GetTag(tag interface{}) (key float64) {
	if node, exists := heap.index[tag]; exists {
		return node.key
	}

	return math.Inf(-1)
}

// GetValue searches and returns the value in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// GetValue will not extract the value so the value will still exist in the heap.
func (heap *FibHeap) GetValue(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
	}

	return
}

// ExtractTag searches and extracts the tag/key in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// ExtractTag will extract the value so the value will no longer exist in the heap.
func (heap *FibHeap) ExtractTag(tag interface{}) (key float64) {
	if node, exists := heap.index[tag]; exists {
		key = node.key
		heap.deleteNode(node)
		return
	}

	return math.Inf(-1)
}

// ExtractValue searches and extracts the value in the heap by the input tag.
// If the input tag does not exist in the heap, nil will be returned.
// ExtractValue will extract the value so the value will no longer exist in the heap.
func (heap *FibHeap) ExtractValue(tag interface{}) (value Value) {
	if node, exists := heap.index[tag]; exists {
		value = node.value
		heap.deleteNode(node)
		return
	}

	return nil
}

// String provides some basic debug information of the heap.
// It returns the total number, roots size, index size and current minimum value of the heap.
// It also returns the topology of the trees by dfs search.
func (heap *FibHeap) String() string {
	var buffer bytes.Buffer

	if heap.num != 0 {
		buffer.WriteString(fmt.Sprintf("Total number: %d, Root Size: %d, Index size: %d,\n", heap.num, heap.roots.Len(), len(heap.index)))
		buffer.WriteString(fmt.Sprintf("Current minimun: key(%f), tag(%v), value(%v),\n", heap.min.key, heap.min.tag, heap.min.value))
		buffer.WriteString(fmt.Sprintf("Heap detail:\n"))
		probeTree(&buffer, heap.roots)
		buffer.WriteString(fmt.Sprintf("\n"))
	} else {
		buffer.WriteString(fmt.Sprintf("Heap is empty.\n"))
	}

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

func (heap *FibHeap) insert(tag interface{}, key float64, value Value) error {
	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage ")
	}

	if _, exists := heap.index[tag]; exists {
		return errors.New("Duplicate tag is not allowed ")
	}

	node := new(node)
	node.children = list.New()
	node.tag = tag
	node.key = key
	node.value = value

	node.self = heap.roots.PushBack(node)
	heap.index[node.tag] = node
	heap.num++

	if heap.min == nil || heap.min.key > node.key {
		heap.min = node
	}

	return nil
}

func (heap *FibHeap) extractMin() *node {
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

	return min
}

func (heap *FibHeap) deleteNode(n *node) {
	heap.decreaseKey(n, n.value, math.Inf(-1))
	heap.ExtractMinValue()
}

func (heap *FibHeap) link(parent, child *node) {
	child.marked = false
	child.parent = parent
	child.self = parent.children.PushBack(child)
	parent.degree++
}

func (heap *FibHeap) resetMin() {
	heap.min = heap.roots.Front().Value.(*node)
	for tree := heap.min.self.Next(); tree != nil; tree = tree.Next() {
		if tree.Value.(*node).key < heap.min.key {
			heap.min = tree.Value.(*node)
		}
	}
}

func (heap *FibHeap) decreaseKey(n *node, value Value, key float64) error {
	if key >= n.key {
		return errors.New("New key is not smaller than current key ")
	}

	n.key = key
	n.value = value
	if n.parent != nil {
		parent := n.parent
		if n.key < n.parent.key {
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}

	if n.parent == nil && n.key < heap.min.key {
		heap.min = n
	}

	return nil
}

func (heap *FibHeap) increaseKey(n *node, value Value, key float64) error {
	if key <= n.key {
		return errors.New("New key is not larger than current key ")
	}

	n.key = key
	n.value = value

	child := n.children.Front()
	for child != nil {
		childNode := child.Value.(*node)
		child = child.Next()
		if childNode.key < n.key {
			heap.cut(childNode)
			heap.cascadingCut(n)
		}
	}

	if heap.min == n {
		heap.resetMin()
	}

	return nil
}

func (heap *FibHeap) cut(n *node) {
	n.parent.children.Remove(n.self)
	n.parent.degree--
	n.parent = nil
	n.marked = false
	n.self = heap.roots.PushBack(n)
}

func (heap *FibHeap) cascadingCut(n *node) {
	if n.parent != nil {
		if !n.marked {
			n.marked = true
		} else {
			parent := n.parent
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}
}
