package main

import (
	"fmt"
	"net/http"

	"github.com/wiretick/go-htmx/components/logger"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

func (Post) create(title, body string) Post {
	return Post{
		Id:    1,
		Title: title,
		Body:  body,
	}
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	post := Post.create("hello", "world")
	fmt.Fprintf(w, "Hello world, from logger: %s", logger.Write("Serving on port 8000"))
	fmt.Fprint(w, "\nAnother line of text")
}

func main() {
	http.HandleFunc("GET /", getPosts)

	http.ListenAndServe(":8000", nil)
}
