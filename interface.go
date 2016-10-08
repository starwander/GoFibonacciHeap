package fibHeap

import (
	"container/list"
)

type FibHeap interface {
	Num() uint
	Insert(Value)
	Minimum() Value
	ExtractMin() interface{}
	Union(FibHeap) FibHeap
	//DecreaseKey()
	//Delete()
}

type Value interface {
	Key() float64
}

func NewFibHeap() FibHeap {
	heap := new(fibHeap)
	heap.roots = list.New()
	heap.num = 0
	heap.min = nil

	return heap
}
