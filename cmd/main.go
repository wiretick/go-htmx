package main

import (
	"log"
	"net/http"

	"github.com/wiretick/go-htmx/core"
	"github.com/wiretick/go-htmx/handlers"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", core.APIHandler(handlers.HandleGetPosts))
	router.HandleFunc("GET /posts/{id}", core.APIHandler(handlers.HandleGetPostByID))
	router.HandleFunc("POST /posts", core.APIHandler(handlers.HandleCreatePost))

	m := core.UseMiddleware(
		core.LoggingMiddleware,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: m(router),
	}

	log.Fatal(server.ListenAndServe())
}
