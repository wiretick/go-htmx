package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	title string
	body  string
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	for _, post := range posts {
		_, err := fmt.Fprintf(w, "%s\n", post.body)
		if err != nil {
			return
		}
	}

	fmt.Fprint(w, "\nAnother line of text\n")
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return
	}

	// super safe way of doing things ;)
	fmt.Fprint(w, posts[id])
}

func createPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	post := Post{title: title, body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
}

var posts []Post = []Post{
	{title: "what", body: "whats"},
	{title: "another", body: "body another"},
}

func main() {
	// Using custom ServeMux makes it possible to have multiple different routers
	router := http.NewServeMux()

	router.HandleFunc("GET /", getPosts)
	router.HandleFunc("GET /posts/{id}", getPost)
	router.HandleFunc("POST /posts", createPost)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
