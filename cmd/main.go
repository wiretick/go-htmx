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
	if err := post.Index(posts).Render(r.Context(), w); err != nil {
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

	if err := post.Show(posts[id]).Render(r.Context(), w); err != nil {
		return err
	}

	return nil
}

func createPost(w http.ResponseWriter, r *http.Request) error {
	title := r.FormValue("title")
	body := r.FormValue("body")

	if len(title) > 20 {
		return InvalidRequest("Title needs to be shorter than 20 characters")
	}

	post := types.Post{Title: title, Body: body}
	posts = append(posts, post)

	w.WriteHeader(http.StatusSeeOther)
	return nil
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
	router.HandleFunc("POST /posts", Make(createPost))

	m := middleware.Use(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8000",
		Handler: m(router),
	}

	log.Fatal(server.ListenAndServe())
}
