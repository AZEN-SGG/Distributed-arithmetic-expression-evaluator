package calculator

import (
	"Distributed-arithmetic-expression-evaluator/rest"
	"context"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	alphabet           = []int32{40, 41, 42, 47, 43, 45}
	canBe              = []int32{48, 49, 50, 51, 52, 43, 45, 42, 40, 41, 47, 53, 54, 55, 56, 57}
	arithmeticExecTime = map[int32]time.Duration{43: time.Millisecond * 500, 45: time.Millisecond * 750,
		42: time.Millisecond * 1000, 47: time.Millisecond * 1500}
)

func Waiter(value1, value2 int, operate int32) int {
	time.Sleep(arithmeticExecTime[operate])

	switch operate {
	case 42:
		return value1 * value2

	case 43:
		return value1 + value2

	case 45:
		return value1 - value2

	case 47:
		return value1 / value2

	default:
		return -1
	}
}

// Delimiter Распределяет подсчёт на разные ярусы
func Delimiter(expression string) ([]int32, []int, error) {
	var queue = make([]int32, 0)
	var values = make([]int, 0)
	var value string

	for _, val := range expression {
		if slices.Contains(alphabet, val) {
			if value != "" {
				num, err := strconv.Atoi(value)

				if err != nil {
					return nil, nil, rest.NewError("Extraneous characters found in expression: %s", value)
				}

				values = append(values, num)
				value = ""
			}

			queue = append(queue, val)
		} else {
			value += string(val)
		}
	}

	if value != "" {
		num, err := strconv.Atoi(value)

		if err != nil {
			return nil, nil, rest.NewError("Extraneous characters found in expression: %s", value)
		}

		values = append(values, num)
	}

	return queue, values, nil
}

// Distributor Разделяет значения на множество действий
func Distributor(queue []int32, values []int) ([][]int32, [][]int, error) {
	var expressions = [][]int32{{}}
	var sortValues = [][]int{{}}

	var indexes = []int{0}
	var indexLastExpression = -1
	var indexValue = 0 // Индекс числа из Value
	var index int

	for i, val := range queue {
		switch val {
		case 40: // (
			indexes = append(indexes, len(expressions))
			sortValues = append(sortValues, make([]int, 0))
			expressions = append(expressions, make([]int32, 0))

		case 41: // )
			index = rest.Last(indexes)

			if indexLastExpression != -1 {
				sortValues[index] = append(sortValues[index], -indexLastExpression)
				indexLastExpression = -1
			} else {
				sortValues[index] = append(sortValues[index], values[indexValue])
				indexValue++
			}

			indexLastExpression = index

			indexes = slices.Delete(indexes, len(indexes)-1, len(indexes))

			if i == len(queue)-1 {
				sortValues[rest.Last(indexes)] = append(sortValues[rest.Last(indexes)], -indexLastExpression)
			}

		default:
			index = rest.Last(indexes)

			if indexLastExpression != -1 {
				sortValues[index] = append(sortValues[index], -indexLastExpression)
				indexLastExpression = -1
			} else {
				sortValues[index] = append(sortValues[index], values[indexValue])
				indexValue++
			}

			expressions[index] = append(expressions[index], val)
		}
	}

	if index == 0 {
		if indexLastExpression != -1 {
			sortValues[index] = append(sortValues[index], -indexLastExpression)
			indexLastExpression = -1
		} else {
			sortValues[index] = append(sortValues[index], values[indexValue])
			indexValue++
		}
	}

	return expressions, sortValues, nil
}

func Mathematician(expressions [][]int32, values [][]int) int {
	var doneCh = make(chan []int)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go Proletarian(&wg, expressions, values, 0, &doneCh)

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	for i := range doneCh {
		return i[0]
	}

	return -1
}

func Proletarian(wg *sync.WaitGroup, expressions [][]int32, values [][]int, index int, outCh *chan []int) {
	defer wg.Done()
	var valueWG = sync.WaitGroup{}

	var expression = expressions[index]
	var value = values[index]
	var calculated = make([]int, len(value))

	var routine = ArithmeticSorter(expression)

	var value1, value2 int

	for _, i := range routine {
		if calculated[i] == 0 {
			value1 = value[i]
		} else {
			value1 = calculated[i]
		}

		if calculated[i+1] == 0 {
			value2 = value[i+1]
		} else {
			value2 = calculated[i+1]
		}

		var valueCh = make(chan []int)

		if value1 <= 0 {
			valueWG.Add(1)
			go Proletarian(&valueWG, expressions, values, -value1, &valueCh)
		}

		if value2 <= 0 {
			valueWG.Add(1)
			go Proletarian(&valueWG, expressions, values, -value2, &valueCh)
		}

		go func() {
			defer close(valueCh)
			valueWG.Wait()
		}()

		for val := range valueCh {
			if val[1] == -value1 {
				value1 = val[0]
			} else {
				value2 = val[0]
			}
		}

		calculated[i] = Waiter(value1, value2, expression[i])
		calculated[i+1] = calculated[i]
	}

	*outCh <- []int{calculated[rest.Last(routine)], index}
}

// ArithmeticSorter Сортирует операции по важности по убыванию
func ArithmeticSorter(expression []int32) []int {
	var routine []int
	var routineStrong []int

	for i, val := range expression {
		if val == 42 || val == 47 {
			routineStrong = append(routineStrong, i)
		} else if val == 43 || val == 45 {
			routine = append(routine, i)
		}
	}

	return append(routineStrong, routine...)
}

// CalculationTime Считает примерное время выполнения операции
func CalculationTime(expression string) (time.Duration, error) {
	expressions, _, err := Delimiter(expression)

	if err != nil {
		return 0, err
	}

	var workingHours time.Duration

	for _, expr := range expressions {
		workingHours += arithmeticExecTime[expr]
	}

	return workingHours, nil
}

// PreparingExpression Проверяет выражение на правильность формулировки и форматирует
func PreparingExpression(expression string) (string, error) {
	var parenthesis = 0

	expression = strings.Replace(expression, " ", "", strings.Count(expression, " "))

	var dataReplace = map[string]string{"+-": "-", "--": "+", "++": "+"}

	for i, val := range expression {
		if val == 40 {
			parenthesis++
		} else if val == 41 {
			parenthesis--
		}

		if parenthesis < 0 {
			return "", rest.NewError("Extra closed parenthesis: %s", strconv.Itoa(i))
		}
	}

	if parenthesis != 0 {
		return "", rest.NewError("Extra open parenthesis")
	}

	expression = strings.Trim(expression, "/*+-")

	for key, val := range dataReplace {
		for strings.Count(expression, key) != 0 {
			expression = strings.Replace(expression, key, val, strings.Count(expression, key))
		}
	}

	for i, elem := range expression {
		if !slices.Contains(canBe, elem) {
			return "", rest.NewError("Foreign character detected: %s", string(elem))
		} else if slices.Contains(alphabet, elem) && elem > 41 && ((slices.Contains(alphabet, int32(expression[i-1])) && expression[i-1] > 41) || (slices.Contains(alphabet, int32(expression[i+1])) && expression[i+1] > 41)) {
			return "", rest.NewError("Incorrect expression: %s", expression[i-1:i+2])
		}
	}

	return expression, nil
}

// Calculator Решает арифметическое выражение
func Calculator(ctx context.Context, expression string) {
	defer close(ctx.Value("errCh").(chan error))
	defer close(ctx.Value("stateCh").(chan int))

	expression, err := PreparingExpression(expression)

	if err != nil {
		ctx.Value("errCh").(chan error) <- err
		return
	}

	queue, values, err := Delimiter(expression)

	if err != nil {
		ctx.Value("errCh").(chan error) <- err
		return
	}

	expressions, value, err := Distributor(queue, values)

	if err != nil {
		ctx.Value("errCh").(chan error) <- err
		return
	}

	for i, val := range expressions {
		if len(val)+1 != len(value[i]) {
			ctx.Value("errCh").(chan error) <- rest.NewError("Incorrect expression")
			return
		} else if len(val) == 0 || len(value[i]) < 2 {
			ctx.Value("errCh").(chan error) <- rest.NewError("Too few arguments")
			return
		}
	}

	answer := Mathematician(expressions, value)
	ctx.Value("stateCh").(chan int) <- answer

	return
}
