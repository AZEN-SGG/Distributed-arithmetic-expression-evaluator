package rest

import (
	"fmt"
	"time"
)

// Пришлось поместить сюда, чтобы не происходит cycle, так как пакет calculator нужен данный класс

type Expression struct {
	Value      int
	Express    string
	Result     chan int
	ErrCh      chan error
	Created    time.Time
	Expiration time.Duration // Время истечения срока
}

func (express *Expression) Close() {
	close(express.ErrCh)
	close(express.Result)
}

func (express *Expression) GetValue() (int, error) {
	select {
	case err := <-express.ErrCh:
		return 0, err
	case answer := <-express.Result:
		if express.Value == -1 {
			express.Value = answer
		}

		return answer, nil
	default:
		return -1, nil
	}
}

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
