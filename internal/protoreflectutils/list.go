package protoreflectutils

import "iter"

type ProtoreflectList[T any] interface {
	Len() int
	Get(i int) T
}

func Iterate[T any](list ProtoreflectList[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < list.Len(); i++ {
			if !yield(list.Get(i)) {
				return
			}
		}
	}
}

func Map[T any, U any](list ProtoreflectList[T], mapFn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for i := 0; i < list.Len(); i++ {
			if !yield(mapFn(list.Get(i))) {
				return
			}
		}
	}
}

func MapToSlice[T any, U any](list ProtoreflectList[T], mapFn func(T) U) []U {
	result := []U{}
	for i := 0; i < list.Len(); i++ {
		result = append(result, mapFn(list.Get(i)))
	}
	return result
}
