package main

import (
	"log"

	"github.com/wiretick/go-htmx/core"
	"github.com/wiretick/go-htmx/handlers"
)

func main() {
	store, err := core.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := handlers.NewAPIServer(":8000", store)
	log.Fatal(server.Run())
}
