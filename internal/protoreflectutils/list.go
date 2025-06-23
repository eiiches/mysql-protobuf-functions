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
