package server

import (
	"Distributed-arithmetic-expression-evaluator/calculator"
	"Distributed-arithmetic-expression-evaluator/data"
	"Distributed-arithmetic-expression-evaluator/expressions"
	"Distributed-arithmetic-expression-evaluator/rest"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	MainExpressions = expressions.NewExpressions()
)

func FormatExpression(id string, expr *rest.Expression) []string {
	var ok, err = expr.GetValue()
	var status string

	switch {
	case err != nil:
		status = err.Error()
	case ok == -1:
		status = "Считается"
	default:
		status = "Высчитан"
	}

	return []string{id, status, expr.Express, expr.Created.Format("02 Jan at 15:04:05"), strconv.FormatInt(expr.Expiration.Milliseconds(), 10) + "ms"}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		return
	}

	var id = r.FormValue("id")
	var result, err = MainExpressions.GetExpression(id)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	if result.Value == -1 {
		result.Value, err = result.GetValue()

		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	_, err = fmt.Fprintf(w, "Expression - %s = %d\nCreation data: %s\nTime: %s", result.Express, result.Value, result.Created, result.Expiration)

	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func ListProcessHandler(w http.ResponseWriter, _ *http.Request) {
	var err error

	_, err = fmt.Fprintln(w, "List of process:")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = fmt.Fprint(w, "Format: ID - state - expression - creation date - approximate calculation time\n")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = fmt.Fprintln(w, "")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	var values = MainExpressions.GetExpressions()

	for id, expr := range values {
		formatExpression := FormatExpression(id, expr)

		_, err = fmt.Fprint(w, strings.Join(formatExpression, " - ")+"\n")

		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
}

func ArithmeticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(400)
		return
	}

	expr := strings.Replace(r.FormValue("expression"), " ", "+", -1)
	id := r.FormValue("id")

	if id != "" {
		var expr, err = calculator.PreparingExpression(expr)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		_, err = MainExpressions.AddExpression(id, expr)

		if err != nil {
			w.WriteHeader(400)
			return
		}
	}
}

func MathOperationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var _, err = fmt.Fprintf(w, "Math operations:\n+ : %s\n- : %s\n* : %s\n/ : %s", calculator.ArithmeticExecTime[43], calculator.ArithmeticExecTime[45], calculator.ArithmeticExecTime[42], calculator.ArithmeticExecTime[47])

		if err != nil {
			w.WriteHeader(500)
			return
		}

	} else if r.Method == http.MethodPost {
		// Обработка POST запроса с обновленными значениями операций
		var err = r.ParseForm()

		if err != nil {
			w.WriteHeader(500)
			return
		}

		addition := r.Form.Get("addition")
		subtraction := r.Form.Get("subtraction")
		multiplication := r.Form.Get("multiplication")
		division := r.Form.Get("division")

		var operations = make([]*calculator.Operation, 0, 4)
		var operation *calculator.Operation

		if addition != "" {
			operation, err = calculator.FormatOperation(43, addition)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			operations = append(operations, operation)
		}

		if subtraction != "" {
			operation, err = calculator.FormatOperation(45, subtraction)

			if err != nil {
				w.WriteHeader(400)
				return
			}

			operations = append(operations, operation)
		}

		if multiplication != "" {
			operation, err = calculator.FormatOperation(42, multiplication)

			if err != nil {
				w.WriteHeader(400)
				return
			}

			operations = append(operations, operation)
		}

		if division != "" {
			operation, err = calculator.FormatOperation(47, division)

			if err != nil {
				w.WriteHeader(400)
				return
			}

			operations = append(operations, operation)
		}

		calculator.MathOperation(operations...)
		err = data.UploadArithmetic(calculator.ArithmeticExecTime, "data/arithmetic.csv")

		if err != nil {
			w.WriteHeader(500)
			return
		}

		// Вывод обновленных значений операций
		_, err = fmt.Fprintf(w, "Operations updated:\n+ : %s\n- : %s\n* : %s\n/ : %s", calculator.ArithmeticExecTime[43], calculator.ArithmeticExecTime[45], calculator.ArithmeticExecTime[42], calculator.ArithmeticExecTime[47])

		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

func ProcessesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var _, err = fmt.Fprint(w, "Processes:\n\n")

		if err != nil {
			w.WriteHeader(500)
			return
		}

		for i, elem := range calculator.ComputingPower {
			_, err = fmt.Fprintf(w, "%d - %s\n", i, string(elem))

			if err != nil {
				w.WriteHeader(500)
				return
			}
		}

	} else {
		w.WriteHeader(400)
		return
	}
}

func MuxHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", ResultHandler)
	mux.HandleFunc("/list", ListProcessHandler)
	mux.HandleFunc("/math", MathOperationsHandler)
	mux.HandleFunc("/processes", ProcessesHandler)
	mux.HandleFunc("/expression", ArithmeticsHandler)

	return mux
}

func StartHandler(port string) {
	var mux = MuxHandler()
	_ = MainExpressions.DownloadExpressions("data/data_expressions.csv")
	_ = data.DownloadArithmetic(calculator.ArithmeticExecTime, "data/arithmetic.csv")

	fmt.Printf("Server start listening on http://localhost:%s/", port)
	var err = http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err.Error())
	}
}
