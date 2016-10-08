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
				demo.tag = i
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
				demo.tag = i
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
			heap.Union(anotherHeap)
			Expect(heap.Minimum()).Should(BeNil())
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the non-empty one into the empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				anotherHeap.Insert(demo)
			}
			number := anotherHeap.Num()
			min := anotherHeap.Minimum()

			heap.Union(anotherHeap)
			Expect(heap.Minimum()).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the empty one into the non-empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}
			number := heap.Num()
			min := heap.Minimum()

			heap.Union(anotherHeap)
			Expect(heap.Minimum()).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given two fibHeap with multiple values, when call ExtractMin api after unioned, it should extract the minimum value inserted into both heaps.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 5000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}
			for i := 0; i < 5000; i++ {
				anotherdemo := new(demoStruct)
				anotherdemo.tag = i + 5000
				anotherdemo.key = rand.Float64()
				anotherdemo.value = fmt.Sprint(anotherdemo.key)
				anotherHeap.Insert(anotherdemo)
			}
			min := heap.Minimum().(*demoStruct).key
			if anotherHeap.Minimum().(*demoStruct).key < min {
				min = anotherHeap.Minimum().(*demoStruct).key
			}
			heap.Union(anotherHeap)

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			lastKey := heap.Minimum().(*demoStruct).key
			Expect(lastKey).Should(BeEquivalentTo(min))
			for i := 0; i < 10000; i++ {
				extracted := heap.ExtractMin().(*demoStruct)
				Expect(extracted.key).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = extracted.key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("index tests", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
			anotherHeap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
			anotherHeap = nil
		})

		It("Given one fibHeap, when Insert values with same tag, it should return an error.", func() {
			demo := new(demoStruct)
			demo.tag = 1
			demo.key = 1
			demo.value = "1"
			err := heap.Insert(demo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(heap.Minimum()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			err = heap.Insert(demo)
			Expect(err).Should(HaveOccurred())
			Expect(heap.Minimum()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given two fibHeaps which both has value with same tag, when call Union, it should return an error.", func() {
			demo := new(demoStruct)
			demo.tag = 1
			demo.key = 1
			demo.value = "1"
			heap.Insert(demo)
			anotherDemo := new(demoStruct)
			anotherDemo.tag = 1
			anotherDemo.key = 2
			anotherDemo.value = "2"
			anotherHeap.Insert(anotherDemo)

			err := heap.Union(anotherHeap)
			Expect(err).Should(HaveOccurred())
			Expect(heap.Minimum()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			Expect(anotherHeap.Minimum()).Should(BeEquivalentTo(anotherDemo))
			Expect(anotherHeap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given one fibHeaps which has not a value with TAG, when GetTag this TAG, it should return nil.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}

			Expect(heap.GetTag(10000)).Should(BeNil())
		})

		It("Given one fibHeaps which has a value with TAG, when GetTag this TAG, it should return the value with TAG.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo)
			}
			tagValue := new(demoStruct)
			tagValue.tag = 10000
			tagValue.key = 10000
			tagValue.value = "10000"
			heap.Insert(tagValue)

			Expect(heap.GetTag(10000)).Should(BeEquivalentTo(tagValue))
			Expect(heap.Num()).Should(BeEquivalentTo(1001))
		})
	})
})

type demoStruct struct {
	tag   int
	key   float64
	value string
}

func (demo *demoStruct) Tag() interface{} {
	return demo.tag
}

func (demo *demoStruct) Key() float64 {
	return demo.key
}
