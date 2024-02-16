package rest

import (
	"fmt"
)

func Last[E any](s []E) E {
	return s[len(s)-1]
}

func NewError(format string, values ...interface{}) error {
	return fmt.Errorf(format, values...)
}

func MapGetKeys[K comparable, V any](m map[K]V) []K {
	var keys = make([]K, len(m))
	var index int

	for key := range m {
		keys[index] = key
		index++
	}

	return keys
}
