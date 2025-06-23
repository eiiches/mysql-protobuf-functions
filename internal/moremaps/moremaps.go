package moremaps

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

func SortedEntries[K cmp.Ordered, V any, Map ~map[K]V](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range slices.Sorted(maps.Keys(m)) {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}
