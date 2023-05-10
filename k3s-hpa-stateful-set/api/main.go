package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func fib(n int) int {
	if n <= 0 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}

func main() {

	// host:3000/workload/:number
	http.HandleFunc("/workload/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		url := r.URL.Path
		segments := strings.Split(url, "/")
		numberstr := segments[len(segments)-1]
		number, _ := strconv.Atoi(numberstr)

		fibonnaci := fmt.Sprintf("%d", fib(number))

		w.Header().Add("x-fib-number", fibonnaci)
	})

	http.ListenAndServe(":3000", nil)
}
