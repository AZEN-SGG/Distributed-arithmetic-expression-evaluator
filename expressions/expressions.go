package expressions

import (
	"Distributed-arithmetic-expression-evaluator/calculator"
	"Distributed-arithmetic-expression-evaluator/rest"
	"slices"
	"sync"
	"time"
)

// Expressions Является структурой выражений
type Expressions struct {
	IDs map[string]*rest.Expression
	mu  sync.Mutex
}

func NewExpressions() *Expressions {
	return &Expressions{
		IDs: map[string]*rest.Expression{},
		mu:  sync.Mutex{},
	}
}

func (express *Expressions) AddExpression(ID, expr string) (string, error) {
	express.mu.Lock()
	var keys = rest.MapGetKeys(express.IDs)
	express.mu.Unlock()

	if slices.Contains(keys, ID) {
		return "", rest.NewError("An expression with ID %s is already exists", ID)
	}

	var ex, err = NewExpression(expr)

	if err != nil {
		return "", err
	}

	express.mu.Lock()
	express.IDs[ID] = ex
	express.mu.Unlock()

	go calculator.Calculator(ex)

	return ID, nil
}

func (express *Expressions) Delete(IDs ...string) {
	defer express.mu.Unlock()
	express.mu.Lock()
	for _, key := range IDs {
		delete(express.IDs, key)
	}
}

func (express *Expressions) Lock() {
	express.mu.Lock()
}

func (express *Expressions) Unlock() {
	express.mu.Unlock()
}

// GetExpression Возвращает результат, -1 или ошибку в зависимости от состояния процесса
func (express *Expressions) GetExpression(ID string) (*rest.Expression, error) {
	express.mu.Lock()
	var expr, ok = express.IDs[ID]
	express.mu.Unlock()

	if !ok {
		return nil, rest.NewError("There is no such expression: %d", ID)
	}

	return expr, nil
}

func (express *Expressions) GetExpressions() map[string]*rest.Expression {
	var expressions = map[string]*rest.Expression{}

	express.mu.Lock()
	for key, value := range express.IDs {
		expressions[key] = value
	}
	express.mu.Unlock()

	return expressions
}

func NewExpression(express string) (*rest.Expression, error) {
	var duration, err = calculator.CalculationTime(express)

	if err != nil {
		return nil, err
	}

	return &rest.Expression{
		Express:    express,
		Result:     make(chan int),
		ErrCh:      make(chan error),
		Created:    time.Now(),
		Expiration: duration,
	}, nil
}
