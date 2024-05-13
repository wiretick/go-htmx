package main

import (
	"fmt"
	"net/http"
)

type Post struct {
	Title string
	Body  string
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	for _, post := range posts {
		_, err := fmt.Fprintf(w, "%s\n", post.Body)
		if err != nil {
			return
		}
	}

	fmt.Fprint(w, "\nAnother line of text\n")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	post := Post{Title: title, Body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
}

var posts []Post = []Post{{Title: "what", Body: "whats"}, {Title: "another", Body: "body another"}}

func main() {
	http.HandleFunc("GET /", getPosts)
	http.HandleFunc("POST /posts", createPost)

	http.ListenAndServe(":8000", nil)
}
