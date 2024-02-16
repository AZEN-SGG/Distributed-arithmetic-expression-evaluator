package expressions

import (
	"Distributed-arithmetic-expression-evaluator/calculator"
	"Distributed-arithmetic-expression-evaluator/rest"
	"context"
	"math/rand/v2"
	"slices"
	"sync"
)

// Expressions Является структурой выражений
type Expressions struct {
	IDs map[int]context.Context
	mu  sync.Mutex
}

func NewExpressions() *Expressions {
	return &Expressions{
		IDs: map[int]context.Context{},
		mu:  sync.Mutex{},
	}
}

func (express *Expressions) AddExpression(expression string) (int, error) {
	var ID = rand.N(999_999_999)

	express.mu.Lock()
	var keys = rest.MapGetKeys(express.IDs)
	express.mu.Unlock()

	for slices.Contains(keys, ID) {
		ID = rand.N(999_999_999)
	}

	var ctx = context.WithValue(context.WithValue(context.WithValue(context.Background(), "errCh", make(chan error)), "stateCh", make(chan int)), "expression", expression)

	express.mu.Lock()
	express.IDs[ID] = ctx
	express.mu.Unlock()

	go calculator.Calculator(ctx, expression)

	return ID, nil
}

func (express *Expressions) Delete(IDs ...int) {
	defer express.mu.Unlock()
	express.mu.Lock()
	for _, key := range IDs {
		delete(express.IDs, key)
	}
}

// GetExpression Возвращает результат, -1 или ошибку в зависимости от состояния процесса
func (express *Expressions) GetExpression(ID int) (int, error) {
	express.mu.Lock()
	var ctx, ok = express.IDs[ID]
	express.mu.Unlock()

	if !ok {
		return 0, rest.NewError("There is no such expression: %d", ID)
	}

	select {
	case err := <-ctx.Value("errCh").(chan error):
		return 0, err
	case answer := <-ctx.Value("stateCh").(chan int):
		return answer, nil
	default:
		return -1, nil
	}
}

func (express *Expressions) GetExpressions() map[int]string {
	var expressions = map[int]string{}

	express.mu.Lock()
	for key, value := range express.IDs {
		expressions[key] = value.Value("expression").(string)
	}
	express.mu.Unlock()

	return expressions
}
