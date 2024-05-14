package main

import (
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
		_, err := w.Write([]byte("\nPost: " + post.body + "\n"))
		if err != nil {
			return
		}
	}

	w.Write([]byte("\nAnother line of text\n"))
}

func getPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Oops, needs to be a valid integer"))
		return
	}

	if id < 0 || id >= len(posts) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Could not find what you are looking for"))
		return
	}

	w.Write([]byte("Amazing content: " + posts[id].body))
}

func createPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	post := Post{title: title, body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
}

var posts []Post = []Post{
	{title: "what", body: "whats long text"},
	{title: "another", body: "body another longer text"},
}

func main() {
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
