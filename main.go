package main

import (
	"net/http"
)

/*
func main() {
	http.HandleFunc("/expression", server.ArithmeticsHandler)
	http.HandleFunc("/list", server.ListProcessHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Start handling")
}

*/

func main() {
	http.HandleFunc("/math-operations", MathOperationsHandler)
	http.ListenAndServe(":8080", nil)
}
