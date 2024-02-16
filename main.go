package main

import (
	"Distributed-arithmetic-expression-evaluator/expressions"
	"fmt"
	"time"
)

func main() {
	var exp = expressions.NewExpressions()
	var ID, _ = exp.AddExpression("(6+6)+6")

	fmt.Println(exp.GetExpression(ID))

	fmt.Println(exp.GetExpressions())
	exp.Delete(ID)
	time.Sleep(time.Second * 10)
	fmt.Println(exp.GetExpression(ID))
}
