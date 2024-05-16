package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/wiretick/go-htmx/types"
	"github.com/wiretick/go-htmx/view/post"
)

// Custom handler signature to return errors
type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	// Need to return a normal HandlerFunc to the router so its happy
	// but first I want to deal with the errors if there are any
	// TODO: want to set the correct status code on the HTTP header
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// TODO: consider trying errors.As() later
			if apiErr, ok := err.(APIError); ok {
				slog.Error("API", "msg", apiErr.Msg)
			} else {
				slog.Error("Internal server error", "err", err.Error())
			}
		}
	}
}

type Middleware func(http.Handler) http.Handler

func Apply(xs ...Middleware) Middleware {
	// using this so its easier to add multiple middlewares to a router
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}

type adaptedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *adaptedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	// Writes all the requests to console
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		adapter := &adaptedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(adapter, r)

		log.Println(adapter.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	// Need this Error() to be compatible with the error interface
	return fmt.Sprintf("api error: %d", e.Status)
}

func InvalidRequest(err string) APIError {
	return APIError{
		Status: http.StatusUnprocessableEntity,
		Msg:    err,
	}
}

func NotFound(err string) APIError {
	return APIError{
		Status: http.StatusUnprocessableEntity,
		Msg:    err,
	}
}

func getPosts(w http.ResponseWriter, r *http.Request) error {
	postPage := types.PostPage{
		Posts: posts,
	}

	if err := post.Index(postPage).Render(r.Context(), w); err != nil {
		return err
	}

	return nil
}

func getPostById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return InvalidRequest("Must provide valid integer for post ID")
	}

	if id < 0 || id >= len(posts) {
		return NotFound("Could not find a post with the given ID")
	}

	if _, err := w.Write([]byte("Amazing content: " + posts.Find(id).Body)); err != nil {
		// TODO: better with a server failed error
		return InvalidRequest("Failed to write content")
	}

	return nil // All good :)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")

	post := types.Post{Title: title, Body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
}

var posts types.Posts = types.Posts{
	{Title: "What an adventure", Body: "whats long text"},
	{Title: "Lorem ipsum title", Body: "body another longer text"},
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", Make(getPosts))
	router.HandleFunc("GET /posts/{id}", Make(getPostById))
	router.HandleFunc("POST /posts", createPost)

	middlewares := Apply(
		Logging,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: middlewares(router),
	}

	log.Fatal(server.ListenAndServe())
}
