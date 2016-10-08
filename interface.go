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
	//DecreaseKey()
	//Delete()
	GetTag(interface{}) Value
	//ExtractTag(interface{}) Value
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
