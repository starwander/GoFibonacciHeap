## GoFibonacciHeap
[![Build Status](https://travis-ci.org/EthanZhuang/GoFibonacciHeap.svg?branch=master)](https://travis-ci.org/EthanZhuang/GoFibonacciHeap)
[![codecov](https://codecov.io/gh/EthanZhuang/GoFibonacciHeap/branch/master/graph/badge.svg)](https://codecov.io/gh/EthanZhuang/GoFibonacciHeap)
[![Go Report Card](https://goreportcard.com/badge/github.com/EthanZhuang/GoFibonacciHeap)](https://goreportcard.com/report/github.com/EthanZhuang/GoFibonacciHeap)
[![GoDoc](https://godoc.org/github.com/EthanZhuang/GoFibonacciHeap?status.svg)](https://godoc.org/github.com/EthanZhuang/GoFibonacciHeap)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

GoFibonacciHeap is a [Golang](https://golang.org/) implementation of [Fibonacci Heap](https://en.wikipedia.org/wiki/Fibonacci_heap).
This implementation is a bit different from the traditional Fibonacci Heap with an index map inside.
Thanks to the index map, the internal struct 'node' no longer need to be exposed outsides the package.
The index map also makes the random access to the values in the heap possible.
But the union operation of this implementation is O(n) rather than O(1) of the traditional implementation.

| Operations                 | Insert | Minimum | ExtractMin | Union | DecreaseKey | Delete    | Get  |
| :------------------------: | :----: | :-----: | :--------: | :---: | :---------: | :-------: | :--: |
| Traditional Implementation | O(1)   | O(1)    | O(log n)¹  | O(1)  | O(1)¹       | O(log n)¹ | N/A  |
| This Implementation        | O(1)   | O(1)    | O(log n)¹  | O(n)  | O(1)¹       | O(log n)¹ | O(1) |
¹ Amortized time.

##Requirements
#####Download this package

    go get github.com/EthanZhuang/GoFibonacciHeap

#####Implements Value interface of this package for all values going to be inserted
```go
// Value is the interface that all values push into or pop from the FibHeap must implement.
type Value interface {
	// Tag returns the unique tag of the value.
	// The tag is used in the index map.
	Tag() interface{}
	// Key returns the key as known as the priority of the value.
	// The valid range of the key is (-inf, +inf].
	Key() float64
}
```
## Supported Operations

* Insert: pushes the input value into the heap.
* Minimum: returns the current minimum value in the heap by key.
* ExtractMin: returns the current minimum value in the heap and then extracts the value from the heap.
* Union: merges the input heap in.
* DecreaseKey: decreases and updates the value in the heap by the input.
* Delete: deletes the value in the heap by the input.
* GetTag: searches and returns the value in the heap by the input tag.
* ExtractTag: searches and extracts the value in the heap by the input tag.
* Num: returns the current total number of values in the heap.
* String: provides some basic debug information of the heap.

## Example

```go
package main

import (
	"fmt"
	"github.com/EthanZhuang/GoFibonacciHeap"
)

type student struct {
	name string
	age  float64
}

func (s *student) Tag() interface{} {
	return s.name
}

func (s *student) Key() float64 {
	return s.age
}

func main() {
	heap := fibHeap.NewFibHeap()

	heap.Insert(&student{"John", 18.3})
	heap.Insert(&student{"Tom", 21.0})
	heap.Insert(&student{"Jessica", 19.4})
	heap.Insert(&student{"Amy", 23.1})

	fmt.Println(heap.Num())     //4
	fmt.Println(heap.Minimum()) //&{John 18.3}
	fmt.Println(heap.Num())     //4

	fmt.Println(heap.ExtractMin()) //&{John 18.3}
	fmt.Println(heap.ExtractMin()) //&{Jessica 19.4}
	fmt.Println(heap.Num())        //2

	amy := heap.GetTag("Amy")
	amy.(*student).age = 16.5
	heap.DecreaseKey(amy)
	fmt.Println(heap.ExtractMin()) //&{Amy 16.5}

	fmt.Println(heap.Num()) //1
	fmt.Println(heap.ExtractTag("Tom")) //&{Tom 21}
	fmt.Println(heap.Num()) //0
}
```

## Reference

[GoDoc](https://godoc.org/github.com/EthanZhuang/GoFibonacciHeap)

## LICENSE

GoFibonacciHeap source code is licensed under the [Apache Licence, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
