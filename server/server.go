package server

import (
	"Distributed-arithmetic-expression-evaluator/calculator"
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

func ListProcessHandler(w http.ResponseWriter, _ *http.Request) {
	var err error

	_, err = fmt.Fprintln(w, "List of process:")

	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = fmt.Fprintln(w, "Format: ID - state - expression - creation date - approximate calculation time")

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
		data := FormatExpression(id, expr)

		_, err = fmt.Fprint(w, strings.Join(data, " - "))

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
		}

		_, err = MainExpressions.AddExpression(id, expr)

		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	w.WriteHeader(200)
	return
}

func MathOperationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, `
   <form method="post">
    <label for="addition">+</label>
    <input type="text" id="addition" name="addition" value="0"><br>
    
    <label for="subtraction">-</label>
    <input type="text" id="subtraction" name="subtraction" value="0"><br>
    
    <label for="multiplication">*</label>
    <input type="text" id="multiplication" name="multiplication" value="0"><br>
    
    <label for="division">/</label>
    <input type="text" id="division" name="division" value="0"><br>
    
    <input type="submit" value="Submit">
   </form>
  `)
	} else if r.Method == http.MethodPost {
		// Обработка POST запроса с обновленными значениями операций
		r.ParseForm()
		addition := r.Form.Get("addition")
		subtraction := r.Form.Get("subtraction")
		multiplication := r.Form.Get("multiplication")
		division := r.Form.Get("division")

		// Вывод обновленных значений операций
		fmt.Fprintf(w, "Updated values:\n+ : %s\n- : %s\n* : %s\n/ : %s", addition, subtraction, multiplication, division)
	}
}
