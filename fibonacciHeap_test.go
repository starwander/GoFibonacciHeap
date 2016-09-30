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
		heap FibHeap
	)

	Context("api test", func() {
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
})

type demoStruct struct {
	key   float64
	value string
}

func (demo *demoStruct) Key() float64 {
	return demo.key
}
