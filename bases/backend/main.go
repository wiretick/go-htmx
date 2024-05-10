package main

import (
	"net/http"
	"fmt"

	"github.com/wiretick/go-htmx/components/logger"
)

func main() {
	logger.Write("Starting the server on port 8000")

	http.HandleFunc("GET /", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	http.ListenAndServe(":8000", nil)
}

