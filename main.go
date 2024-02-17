package main

import (
	"Distributed-arithmetic-expression-evaluator/expressions"
	"Distributed-arithmetic-expression-evaluator/server"
	"fmt"
	"net/http"
)

func main() {
	var exp = expressions.NewExpressions()
	var ID, _ = exp.AddExpression("(6+6)+6")

	fmt.Println(exp.GetExpression(ID))

	fmt.Println(exp.GetExpressions())
	fmt.Println(exp.GetExpression(ID))

	http.HandleFunc("/", server.ListProcess)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Start handling")
}
