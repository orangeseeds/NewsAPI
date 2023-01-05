package collectutil

import (
	"sort"
)

type ordered interface {
	~int | ~float32 | ~float64 | ~int32 | ~int64 | ~string
}

func min[T ordered](i, j T) T {
	if i <= j {
		return i
	}
	return j
}
func max[T ordered](i, j T) T {
	if i >= j {
		return i
	}
	return j

}

func Chunk[E any](arr []E, size int) [][]E {
	result := [][]E{}
	for i := 0; i < len(arr); i = i + size {
		start := i
		end := min(start+size, len(arr))
		result = append(result, arr[start:end])
	}
	return result
}

func Map[E, V any](arr []E, callback func(E) V) *[]V {
	result := []V{}
	for _, item := range arr {
		result = append(result, callback(item))
	}
	return &result
}

func Filter[E any](arr []E, callback func(E) bool) *[]E {
	result := []E{}
	for _, item := range arr {
		if callback(item) {
			result = append(result, item)
		}
	}
	return &result
}

func Sort[E any](arr []E, callback func(i, j int) bool) {
	sort.Slice(arr, callback)
}

func Reduce[E, V any](arr []E, callback func(carry V, obj E) V) V {
	var carry V
	for _, item := range arr {
		carry = callback(carry, item)
	}
	return carry
}

func Find[E any](arr []E, callback func(E) bool) *E {
	for _, item := range arr {
		found := callback(item)
		if found {
			return &item
		}
	}
	return nil
}

func Last[E any](arr []E) *E {
	arrLen := len(arr)
	if arrLen == 0 {
		return nil
	}
	return &arr[arrLen-1]
}

func GetSlice[E any](arr []E, base int, size int) []E {
	arrLen := len(arr)
	start := min(base, arrLen)
	end := min(size, arrLen)
	if start >= end {
		return []E{}
	}
	return arr[start:end]
}

// func sortSlice[T Ordered](s []T, ascending bool) {
// 	sort.Slice(s, func(i, j int) bool {
// 		if ascending {
// 			return s[i] < s[j]
// 		}
// 		return !(s[i] < s[j])
// 	})
// }
// type Wrapper[E, V any] interface {
// 	Map(func(E) V) *[]V
// 	Filter(func(E) bool) *[]E
// 	Sort(func(int, int) bool)
// 	Reduce(func(V, E) V) V
// }

// type wrapper[E, V any] []E

// func Wrap[E, V any](arr []E) Wrapper[E, V] {
// 	return wrapper[E, V](arr)
// }

// func (w wrapper[E, V]) Map(callback func(E) V) *[]V {
// 	result := []V{}
// 	for _, item := range w {
// 		result = append(result, callback(item))
// 	}

// 	return &result
// }

// func (w wrapper[E, _]) Filter(callback func(E) bool) *[]E {
// 	result := []E{}
// 	for _, item := range w {
// 		if callback(item) {
// 			result = append(result, item)
// 		}
// 	}
// 	return &result
// }

// func (w wrapper[E, _]) Sort(callback func(i, j int) bool) {
// 	sort.Slice(w, callback)
// }

// func (w wrapper[E, V]) Reduce(callback func(x V, y E) V) V {
// 	var carry V
// 	for _, item := range w {
// 		carry = callback(carry, item)
// 	}
// 	return carry
// }

// type Collection[T any] interface {
// 	All() *[]T
// 	Push(T)
// 	Concat([]T)
// 	Count() int
// 	Filter(func(T) bool) *Collection[T]
// 	First() *T
// 	FirstOrFail() (*T, error)
// 	Map(func(T) T) *Collection[T]
// 	// Search()
// 	Slice(int, int) *[]T
// 	// SortBy()
// 	// Unique(UniqueCallback[T]) *Collection[T]
// 	// Pluck()
// }

// type collection[T any] struct {
// 	data []T
// }

// func NewCollection[T any]() Collection[T] {
// 	return &collection[T]{
// 		data: []T{},
// 	}
// }

// func (c *collection[T]) All() *[]T {
// 	duplicate := c.data
// 	return &duplicate
// }

// func (c *collection[T]) Push(entry T) {
// 	c.data = append(c.data, entry)
// }

// func (c *collection[T]) Concat(array []T) {
// 	c.data = append(c.data, array...)
// }

// func (c *collection[T]) First() *T {
// 	if len(c.data) < 1 {
// 		return nil
// 	}
// 	val := c.data[0]
// 	return &val
// }

// func (c *collection[T]) FirstOrFail() (*T, error) {
// 	firstVal := c.First()
// 	if firstVal == nil {
// 		return nil, errors.New("no elements in collection")
// 	}
// 	return firstVal, nil
// }

// func (c *collection[T]) Count() int {
// 	return len(c.data)
// }

// func (c *collection[T]) Slice(skip int, limit int) *[]T {
// 	start := int(math.Min(float64(skip), float64(len(c.data))))
// 	end := int(math.Min(float64(limit), float64(len(c.data))))
// 	if start >= end {
// 		return &[]T{}
// 	}
// 	data := c.data[start:end]
// 	return &data
// }

// func (c *collection[T]) Filter(callback func(value T) bool) *Collection[T] {
// 	filteredCollection := NewCollection[T]()
// 	for _, val := range c.data {
// 		addToCollection := callback(val)
// 		if addToCollection {
// 			filteredCollection.Push(val)
// 		}
// 	}
// 	return &filteredCollection
// }

// func (c *collection[T]) Map(callback func(value T) T) *Collection[T] {
// 	mappedCollection := NewCollection[T]()
// 	for _, val := range c.data {
// 		changedValue := callback(val)
// 		mappedCollection.Push(changedValue)
// 	}
// 	return &mappedCollection
// }

// // func (c *collection[T]) Unique(callback func(value T) interface{}) *Collection[T] {
// // 	uniqueCollection := NewCollection[T]()
// // 	uniqueMap := map[interface{}]T{}
// // 	for _, val := range c.data {

// // 		attribute := callback(val)
// // 		uniqueMap[attribute] = val
// // 	}
// // 	for _, val := range uniqueMap {
// // 		uniqueCollection.Push(val)
// // 	}
// // 	return &uniqueCollection
// // }

// func Unique[T any, Y comparable](c Collection[T], callback func(value T) (attribute Y)) *Collection[T] {
// 	uniqueCollection := NewCollection[T]()
// 	unique := []Y{}

// 	collectionEntries := *c.All()
// 	for i, val := range collectionEntries {
// 		attr := callback(val)
// 		if slices.Contains(unique, attr) {
// 			collectionEntries = slices.Delete(collectionEntries, i, i+1)
// 		}
// 		unique = append(unique, attr)
// 	}
// 	uniqueCollection.Concat(collectionEntries)
// 	return &uniqueCollection
// }

// type Ordered interface {
// 	~int | ~float32 | ~float64 | ~string
// }

// func sortSlice[T Ordered](s []T, ascending bool) {
// 	sort.Slice(s, func(i, j int) bool {
// 		if ascending {
// 			return s[i] < s[j]
// 		}
// 		return !(s[i] < s[j])
// 	})
// }

// func SortBy[T any, Y Ordered](c Collection[T], callback func(value T) (attribute Y), ascending bool) *Collection[T] {
// 	sortedCollection := NewCollection[T]()
// 	valMap := map[Y]T{}
// 	mapKeys := []Y{}

// 	collectionEntries := *c.All()
// 	for _, val := range collectionEntries {
// 		attr := callback(val)
// 		mapKeys = append(mapKeys, attr)
// 		valMap[attr] = val
// 	}

// 	sortSlice(mapKeys, ascending)
// 	for _, v := range mapKeys {
// 		sortedCollection.Push(valMap[v])
// 	}
// 	return &sortedCollection
// }

// func Collect[T any](array []T) *collection[T] {
// 	return &collection[T]{
// 		data: array,
// 	}
// }
