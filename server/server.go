package server

import (
	"Distributed-arithmetic-expression-evaluator/expressions"
	"Distributed-arithmetic-expression-evaluator/rest"
	"fmt"
	"maps"
	"net/http"
)

var MainExpressions = expressions.NewExpressions()

func FormatExpression(expr *rest.Expression) []string {
	var ok, err = expr.GetValue()
	var status string

	switch {
	case err != nil:
		status = err.Error()
	case ok == 0:
		status = "Высчитан"
	case ok == -1:
		status = "Считается"
	}

	return []string{status, expr.Express, expr.Created.Format("15:04:05 02 Jan")}
}

func ListProcess(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of process:")

	var values = make(map[int]*rest.Expression)

	MainExpressions.Lock()
	maps.Copy(MainExpressions.GetExpressions(), values)
	MainExpressions.Unlock()

	for _, expr := range values {
		fmt.Fprintln(w, FormatExpression(expr))
	}
}
