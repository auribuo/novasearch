package types

import "fmt"

type Tuple[K, V any] struct {
	Key   K
	Value V
}

func NewTuple[K, V any](key K, value V) Tuple[K, V] {
	return Tuple[K, V]{
		Key:   key,
		Value: value,
	}
}

func (tuple Tuple[K, V]) String() string {
	return fmt.Sprintf("(%v, %v)", tuple.Key, tuple.Value)
}
