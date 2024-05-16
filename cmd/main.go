package main

import (
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/wiretick/go-htmx/core"
	"github.com/wiretick/go-htmx/types"
	"github.com/wiretick/go-htmx/view/post"
)

// Custom handler signature to return errors
type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	// The Make function is an adapter between the default http.HandlerFunc
	// and the custom handler APIFunc which returns an error. Because I want
	// to be able to centralize the error handling
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// TODO: consider trying errors.As() later
			if apiErr, ok := err.(APIError); ok {
				w.WriteHeader(apiErr.Status)
				w.Write([]byte("Error: " + apiErr.Msg))
				slog.Error("API", "detail", apiErr.Msg)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Internal server error", "err", err.Error())
			}
		}
	}
}

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	// Need this Error() to make APIError compatible with the error interface
	return ""
}

func InvalidRequest(err string) APIError {
	return APIError{
		Status: http.StatusUnprocessableEntity,
		Msg:    err,
	}
}

func NotFound(err string) APIError {
	return APIError{
		Status: http.StatusNotFound,
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

	return nil
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
	{Title: "We need more", Body: "Need to move on to the next topic now"},
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", Make(getPosts))
	router.HandleFunc("GET /posts/{id}", Make(getPostById))
	router.HandleFunc("POST /posts", createPost)

	m := middleware.Use(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: m(router),
	}

	log.Fatal(server.ListenAndServe())
}
