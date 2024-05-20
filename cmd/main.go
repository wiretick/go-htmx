package main

import (
	"log"

	"github.com/wiretick/go-htmx/core"
)

func main() {
	store, err := core.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8000", store)
	server.Run()
}
