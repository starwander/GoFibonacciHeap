package fibHeap

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math"
	"math/rand"
	"time"
)

var _ = Describe("Tests of fibHeap", func() {
	var (
		heap        FibHeap
		anotherHeap FibHeap
	)

	Context("behaviour tests", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given an empty fibHeap, when call Minimum api, it should return nil.", func() {
			Expect(heap.Minimum()).Should(BeNil())
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call Minimum api, it should return the minimum value inserted.", func() {
			min := math.Inf(1)
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				demo := new(demoStruct)
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
				if demo.key < min {
					min = demo.key
				}
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			Expect(heap.Minimum().(*demoStruct).key).Should(Equal(min))
			Expect(heap.Minimum().(*demoStruct).value).Should(Equal(fmt.Sprint(min)))
			Expect(heap.Num()).Should(BeEquivalentTo(10000))
		})

		It("Given an empty fibHeap, when call ExtractMin api, it should return nil.", func() {
			Expect(heap.ExtractMin()).Should(BeNil())
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call ExtractMin api, it should extract the minimum value inserted.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				demo := new(demoStruct)
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			lastKey := heap.Minimum().(*demoStruct).key
			for i := 0; i < 10000; i++ {
				extracted := heap.ExtractMin().(*demoStruct)
				Expect(extracted.key).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = extracted.key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("union tests", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
			anotherHeap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
			anotherHeap = nil
		})

		It("Given two empty fibHeaps, when call Union api, it should return an empty fibHeap.", func() {
			newHeap := heap.Union(anotherHeap)
			Expect(newHeap.Minimum()).Should(BeNil())
			Expect(newHeap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the non-empty one into the empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				anotherHeap.Insert(demo)
			}
			number := anotherHeap.Num()
			min := anotherHeap.Minimum()

			newHeap := heap.Union(anotherHeap)
			Expect(newHeap.Minimum()).Should(BeEquivalentTo(min))
			Expect(newHeap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the empty one into the non-empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}
			number := heap.Num()
			min := heap.Minimum()

			newHeap := heap.Union(anotherHeap)
			Expect(newHeap.Minimum()).Should(BeEquivalentTo(min))
			Expect(newHeap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given two fibHeap with multiple values, when call ExtractMin api after unioned, it should extract the minimum value inserted into both heaps.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 5000; i++ {
				demo := new(demoStruct)
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}
			for i := 0; i < 5000; i++ {
				anotherdemo := new(demoStruct)
				anotherdemo.key = rand.Float64()
				anotherdemo.value = fmt.Sprint(anotherdemo.key)
				anotherHeap.Insert(anotherdemo)
			}
			min := heap.Minimum().(*demoStruct).key
			if anotherHeap.Minimum().(*demoStruct).key < min {
				min = anotherHeap.Minimum().(*demoStruct).key
			}
			newHeap := heap.Union(anotherHeap)

			Expect(newHeap.Num()).Should(BeEquivalentTo(10000))
			lastKey := newHeap.Minimum().(*demoStruct).key
			Expect(lastKey).Should(BeEquivalentTo(min))
			for i := 0; i < 10000; i++ {
				extracted := newHeap.ExtractMin().(*demoStruct)
				Expect(extracted.key).Should(BeNumerically(">=", lastKey))
				Expect(newHeap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = extracted.key
			}
			Expect(newHeap.Num()).Should(BeEquivalentTo(0))
		})
	})
})

type demoStruct struct {
	key   float64
	value string
}

func (demo *demoStruct) Key() float64 {
	return demo.key
}
