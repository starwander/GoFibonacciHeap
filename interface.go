package fibHeap

import (
	"container/list"
)

type FibHeap interface {
	Num() uint
	Insert(Value) error
	Minimum() Value
	ExtractMin() Value
	Union(FibHeap) error
	DecreaseKey(Value) error
	Delete(Value) error
	GetTag(interface{}) Value
	ExtractTag(interface{}) Value
	String() string
}

type Value interface {
	Tag() interface{}
	Key() float64
}

func NewFibHeap() FibHeap {
	heap := new(fibHeap)
	heap.roots = list.New()
	heap.index = make(map[interface{}]*list.Element)
	heap.num = 0
	heap.min = nil

	return heap
}
