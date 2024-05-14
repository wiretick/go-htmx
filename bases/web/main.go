package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

type Post struct {
	title string
	body  string
}

type Posts []Post

func (p Posts) find(id int) Post {
	// pretty useless, but can see how this could be a search function instead
	return p[id]
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

	w.Write([]byte("Amazing content: " + posts.find(id).body))
}

func createPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	post := Post{title: title, body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
}

var posts Posts = Posts{
	{title: "what", body: "whats long text"},
	{title: "another", body: "body another longer text"},
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", getPosts)
	router.HandleFunc("GET /posts/{id}", getPost)
	router.HandleFunc("POST /posts", createPost)

	//v1 := http.NewServeMux()
	//v1.Handle("/v1/", http.StripPrefix("/v1", router))

	middlewares := CreateStack(
		Logging,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: middlewares(router),
	}

	log.Fatal(server.ListenAndServe())
}
