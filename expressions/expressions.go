package expressions

import (
	"Distributed-arithmetic-expression-evaluator/calculator"
	"Distributed-arithmetic-expression-evaluator/data"
	"Distributed-arithmetic-expression-evaluator/rest"
	"slices"
	"strconv"
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

	err = express.UploadExpressions("data/data_expressions.csv")

	if err != nil {
		return "", err
	}

	return ID, nil
}

func (express *Expressions) Delete(IDs ...string) {
	defer express.mu.Unlock()
	express.mu.Lock()
	for _, key := range IDs {
		delete(express.IDs, key)
	}

	_ = express.UploadExpressions("data/data_expressions.csv")
}

func (express *Expressions) Lock() {
	express.mu.Lock()
}

func (express *Expressions) Unlock() {
	express.mu.Unlock()
}

func (express *Expressions) DownloadExpressions(name string) error {
	var info, err = data.OpenCSV(name, ';')

	if err != nil {
		return err
	}

	for i, val := range info {
		if i == 0 {
			continue
		}

		if err != nil {
			return err
		}

		if val[2] != "-1" {
			var expr, err = NewExpression(val[1])

			digit, err := strconv.Atoi(val[2])

			if err != nil {
				return err
			}

			expr.Value = digit

			express.Lock()
			express.IDs[val[0]] = expr
			express.Unlock()
		} else {
			_, err = express.AddExpression(val[0], val[1])

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (express *Expressions) UploadExpressions(name string) error {
	var csvFile = make([][]string, 0)

	csvFile = append(csvFile, []string{"ID", "Expression", "Value"})

	for key, val := range express.GetExpressions() {
		var expr = []string{key, val.Express, strconv.Itoa(val.Value)}

		csvFile = append(csvFile, expr)
	}

	var err = data.WriteCSV(csvFile, name, ';')
	return err
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
		Value:      -1,
		Express:    express,
		Result:     make(chan int),
		ErrCh:      make(chan error),
		Created:    time.Now(),
		Expiration: duration,
	}, nil
}
