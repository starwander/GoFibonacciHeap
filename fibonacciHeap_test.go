// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

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
		heap        *FibHeap
		anotherHeap *FibHeap
	)

	Context("behaviour tests of tag/key interfaces", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given an empty fibHeap, when call Minimum api, it should return nil.", func() {
			tag, key := heap.Minimum()
			Expect(tag).Should(BeNil())
			Expect(key).Should(BeEquivalentTo(math.Inf(-1)))
		})

		It("Given a empty fibHeap, when call Insert api with a nil value, it should return error.", func() {
			Expect(heap.Insert(nil, 0.0)).Should(HaveOccurred())
		})

		It("Given a empty fibHeap, when call Insert api with a negetive infinity key, it should return error.", func() {
			Expect(heap.Insert(1000, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Minimum api, it should return the minimum value inserted.", func() {
			min := math.Inf(1)
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				key := rand.Float64()
				heap.Insert(i, key)
				if key < min {
					min = key
				}
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(10000))
		})

		It("Given an empty fibHeap, when call ExtractMin api, it should return nil.", func() {
			tag, _ := heap.ExtractMin()
			Expect(tag).Should(BeNil())
		})

		It("Given a fibHeap inserted multiple values, when call ExtractMin api, it should extract the minimum value inserted.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				key := rand.Float64()
				heap.Insert(i, key)
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, lastKey := heap.Minimum()
			for i := 0; i < 10000; i++ {
				_, key := heap.ExtractMin()
				Expect(key).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap, when call DecreaseKey api with a nil value, it should return error.", func() {
			Expect(heap.DecreaseKey(nil, 0.0)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreaseKey(1000, float64(999))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call DecreaseKey api with a negetive infinity key, it should return error.", func() {
			heap.Insert(1000, float64(1000))
			Expect(heap.DecreaseKey(1000, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a greater key, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreaseKey(999, float64(1000))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a smaller key, it should decrease the key of the value in the heap.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i+1000))
			}
			heap.ExtractMinValue()
			for i := 999; i >= 1; i-- {
				Expect(heap.DecreaseKey(i, float64(i))).ShouldNot(HaveOccurred())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(999))

			for i := 1; i < 1000; i++ {
				tag, key := heap.ExtractMin()
				Expect(tag).Should(BeEquivalentTo(i))
				Expect(key).Should(BeEquivalentTo(i))
			}
		})

		It("Given a fibHeap, when call Delete api with a nil value, it should return error.", func() {
			Expect(heap.Delete(nil)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Delete api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.Delete(10000)).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call Delete api, it should remove the value from the heap.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			for i := 0; i < 1000; i++ {
				Expect(heap.Delete(i)).ShouldNot(HaveOccurred())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("behaviour tests of value interfaces", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given an empty fibHeap, when call Minimum api, it should return nil.", func() {
			Expect(heap.MinimumValue()).Should(BeNil())
		})

		It("Given a empty fibHeap, when call Insert api with a nil value, it should return error.", func() {
			Expect(heap.InsertValue(nil)).Should(HaveOccurred())
		})

		It("Given a empty fibHeap, when call Insert api with a negetive infinity key, it should return error.", func() {
			demo := new(demoStruct)
			demo.tag = 1000
			demo.key = math.Inf(-1)
			demo.value = fmt.Sprint(1000)

			Expect(heap.InsertValue(demo)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Minimum api, it should return the minimum value inserted.", func() {
			min := math.Inf(1)
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
				if demo.key < min {
					min = demo.key
				}
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			Expect(heap.MinimumValue().(*demoStruct).key).Should(Equal(min))
			Expect(heap.MinimumValue().(*demoStruct).value).Should(Equal(fmt.Sprint(min)))
			Expect(heap.Num()).Should(BeEquivalentTo(10000))
		})

		It("Given an empty fibHeap, when call ExtractMin api, it should return nil.", func() {
			Expect(heap.ExtractMinValue()).Should(BeNil())
		})

		It("Given a fibHeap inserted multiple values, when call ExtractMin api, it should extract the minimum value inserted.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 10000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			lastKey := heap.MinimumValue().(*demoStruct).key
			for i := 0; i < 10000; i++ {
				extracted := heap.ExtractMinValue().(*demoStruct)
				Expect(extracted.key).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = extracted.key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap, when call DecreaseKey api with a nil value, it should return error.", func() {
			Expect(heap.DecreaseKeyValue(nil)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}

			decreaseDemo := new(demoStruct)
			decreaseDemo.tag = 1000
			decreaseDemo.key = float64(999)
			decreaseDemo.value = fmt.Sprint(decreaseDemo.key)

			Expect(heap.DecreaseKeyValue(decreaseDemo)).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call DecreaseKey api with a negetive infinity key, it should return error.", func() {
			demo := new(demoStruct)
			demo.tag = 1000
			demo.key = float64(1000)
			demo.value = fmt.Sprint(1000)
			heap.InsertValue(demo)

			demo.key = math.Inf(-1)
			Expect(heap.DecreaseKeyValue(demo)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a greater key, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}

			decreaseDemo := new(demoStruct)
			decreaseDemo.tag = 999
			decreaseDemo.key = float64(1000)
			decreaseDemo.value = fmt.Sprint(decreaseDemo.key)

			Expect(heap.DecreaseKeyValue(decreaseDemo)).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a smaller key, it should decrease the key of the value in the heap.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i + 1000)
				demo.value = fmt.Sprint(i)
				heap.InsertValue(demo)
			}
			heap.ExtractMinValue()
			for i := 999; i >= 1; i-- {
				decreaseDemo := new(demoStruct)
				decreaseDemo.tag = i
				decreaseDemo.key = float64(i)
				decreaseDemo.value = fmt.Sprint(i)
				Expect(heap.DecreaseKeyValue(decreaseDemo)).ShouldNot(HaveOccurred())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(999))

			for i := 1; i < 1000; i++ {
				value := heap.ExtractMinValue()
				Expect(value.(*demoStruct).tag).Should(BeEquivalentTo(i))
				Expect(value.(*demoStruct).key).Should(BeEquivalentTo(i))
				Expect(value.(*demoStruct).value).Should(BeEquivalentTo(fmt.Sprint(i)))
			}
		})

		It("Given a fibHeap, when call Delete api with a nil value, it should return error.", func() {
			Expect(heap.DeleteValue(nil)).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Delete api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(i)
				heap.InsertValue(demo)
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			deleteDemo := new(demoStruct)
			deleteDemo.tag = 10000
			deleteDemo.key = float64(10000)
			deleteDemo.value = fmt.Sprint(10000)
			Expect(heap.DeleteValue(deleteDemo)).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call Delete api, it should remove the value from the heap.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(i)
				heap.InsertValue(demo)
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			for i := 0; i < 1000; i++ {
				deleteDemo := new(demoStruct)
				deleteDemo.tag = i
				deleteDemo.key = float64(i)
				deleteDemo.value = fmt.Sprint(i)
				Expect(heap.DeleteValue(deleteDemo)).ShouldNot(HaveOccurred())
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
			Expect(heap.MinimumValue()).Should(BeNil())
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the non-empty one into the empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				anotherHeap.InsertValue(demo)
			}
			number := anotherHeap.Num()
			min := anotherHeap.MinimumValue()

			heap.Union(anotherHeap)
			Expect(heap.MinimumValue()).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the empty one into the non-empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}
			number := heap.Num()
			min := heap.MinimumValue()

			heap.Union(anotherHeap)
			Expect(heap.MinimumValue()).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given two fibHeap with multiple values, when call ExtractMin api after unioned, it should extract the minimum value inserted into both heaps.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 5000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}
			for i := 0; i < 5000; i++ {
				anotherdemo := new(demoStruct)
				anotherdemo.tag = i + 5000
				anotherdemo.key = rand.Float64()
				anotherdemo.value = fmt.Sprint(anotherdemo.key)
				anotherHeap.InsertValue(anotherdemo)
			}
			min := heap.MinimumValue().(*demoStruct).key
			if anotherHeap.MinimumValue().(*demoStruct).key < min {
				min = anotherHeap.MinimumValue().(*demoStruct).key
			}
			heap.Union(anotherHeap)

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			lastKey := heap.MinimumValue().(*demoStruct).key
			Expect(lastKey).Should(BeEquivalentTo(min))
			for i := 0; i < 10000; i++ {
				extracted := heap.ExtractMinValue().(*demoStruct)
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
			err := heap.InsertValue(demo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(heap.MinimumValue()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			err = heap.InsertValue(demo)
			Expect(err).Should(HaveOccurred())
			Expect(heap.MinimumValue()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given two fibHeaps which both has value with same tag, when call Union, it should return an error.", func() {
			demo := new(demoStruct)
			demo.tag = 1
			demo.key = 1
			demo.value = "1"
			heap.InsertValue(demo)
			anotherDemo := new(demoStruct)
			anotherDemo.tag = 1
			anotherDemo.key = 2
			anotherDemo.value = "2"
			anotherHeap.InsertValue(anotherDemo)

			err := heap.Union(anotherHeap)
			Expect(err).Should(HaveOccurred())
			Expect(heap.MinimumValue()).Should(BeEquivalentTo(demo))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			Expect(anotherHeap.MinimumValue()).Should(BeEquivalentTo(anotherDemo))
			Expect(anotherHeap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given one fibHeaps which has not a value with TAG, when GetTag this TAG, it should return nil.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}

			Expect(heap.GetValue(10000)).Should(BeNil())
		})

		It("Given one fibHeaps which has a value with TAG, when GetTag this TAG, it should return the value with this TAG.", func() {
			rand.Seed(time.Now().Unix())
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}
			tagValue := new(demoStruct)
			tagValue.tag = 10000
			tagValue.key = 10000
			tagValue.value = "10000"
			heap.InsertValue(tagValue)

			Expect(heap.GetValue(10000)).Should(BeEquivalentTo(tagValue))
			Expect(heap.Num()).Should(BeEquivalentTo(1001))
		})

		It("Given one fibHeaps which has not a value with TAG, when ExtractTag this TAG, it should return nil.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractValue(1000)).Should(BeNil())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given one fibHeaps which has a value with TAG, when ExtractTag this TAG, it should extract the value with this TAG from the heap.", func() {
			for i := 0; i < 1000; i++ {
				demo := new(demoStruct)
				demo.tag = i
				demo.key = float64(i)
				demo.value = fmt.Sprint(demo.key)
				heap.InsertValue(demo)
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractValue(999).(*demoStruct).value).Should(BeEquivalentTo(fmt.Sprint(999)))
			Expect(heap.Num()).Should(BeEquivalentTo(999))
			Expect(heap.MinimumValue().(*demoStruct).value).Should(BeEquivalentTo(fmt.Sprint(0)))
		})
	})

	Context("debug test", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given one fibHeaps which some values, when call String api, it should retern the internal debug string.", func() {
			for i := 1; i < 5; i++ {
				for j := 10; j < 15; j++ {
					demo := new(demoStruct)
					demo.tag = i * j
					demo.key = float64(i * j)
					demo.value = fmt.Sprint(demo.key)
					heap.InsertValue(demo)
				}
				heap.ExtractMinValue()
			}

			debugMsg := "Total number: 16, Root Size: 1, Index size: 16,\n" +
				"Current minimun: key(14.000000), tag(14), value(&{14 14 14}),\n" +
				"Heap detail:\n" +
				"< 14.000000 < 56.000000 28.000000 < 42.000000 > 30.000000 < 33.000000 36.000000 < 39.000000 > > 20.000000 < 22.000000 24.000000 < 26.000000 > 40.000000 < 44.000000 48.000000 < 52.000000 > > > > > \n"
			Expect(heap.String()).Should(BeEquivalentTo(debugMsg))
		})
	})

	Context("benchmark", func() {
		BeforeEach(func() {
			heap = NewFibHeap()
		})

		AfterEach(func() {
			heap = nil
		})

		Measure("Benchmark Go Fibonacci Heap", func(b Benchmarker) {
			rand.Seed(time.Now().Unix())
			b.Time("1000000 radom operations", func() {
				var (
					insert, minimun, extract, decrease, get, delete int64
					min                                             *demoStruct
				)
				for i := 0; i < 1000000; i++ {
					if i%3 == 0 {
						demo := new(demoStruct)
						demo.tag = i
						demo.key = rand.Float64()
						demo.value = fmt.Sprint(demo.key)
						Expect(heap.InsertValue(demo)).ShouldNot(HaveOccurred())
						insert++
						if min == nil || demo.key < min.key {
							min = demo
						}
					}
					if i%5 == 0 {
						if extracted := heap.ExtractMinValue(); extracted != nil {
							extract++
							Expect(extracted.(*demoStruct).key).Should(BeEquivalentTo(min.key))
							if currentMin := heap.MinimumValue(); currentMin != nil {
								minimun++
								min = currentMin.(*demoStruct)
							} else {
								Expect(heap.Num()).Should(BeEquivalentTo(0))
								min = nil
							}
						}
					}
					if i%7 == 0 {
						if currentMin := heap.MinimumValue(); currentMin != nil {
							if min != nil {
								minimun++
								Expect(currentMin.(*demoStruct).key).Should(BeEquivalentTo(min.key))
							}
						}
					}
					if i%11 == 0 {
						if temp := heap.GetValue(int(3 * rand.Int31n(int32(i/3)+1))); temp != nil {
							get++
							temp.(*demoStruct).key = temp.(*demoStruct).key / 2
							heap.DecreaseKeyValue(temp)
							decrease++
							currentMin := heap.MinimumValue()
							Expect(currentMin).ShouldNot(BeNil())
							minimun++
							min = currentMin.(*demoStruct)
						}
					}
					if i%13 == 0 {
						if temp := heap.GetValue(int(3 * rand.Int31n(int32(i/3)+1))); temp != nil {
							get++
							heap.DeleteValue(temp)
							delete++
							if min != nil && temp.Tag() == min.tag {
								if currentMin := heap.MinimumValue(); currentMin != nil {
									minimun++
									min = currentMin.(*demoStruct)
								} else {
									Expect(heap.Num()).Should(BeEquivalentTo(0))
									min = nil
								}
							}
						}
					}
				}
				fmt.Println("Final heap size:", heap.Num())
				fmt.Println("Total insert:", insert, "Total minimun:", minimun, "Total extract:", extract, "Total get:", get, "Total decrease:", decrease, "Total delete:", delete)
				Expect(heap.Num()).Should(BeEquivalentTo(insert - extract - delete))
			})
		}, 10)
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
