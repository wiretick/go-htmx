package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ThirdPartyResponse struct {
	val string
	err error
}

func thirdPartyApiCall() (string, error) {
	time.Sleep(time.Millisecond * 500)
	return "data retrieved", nil
}

func (s *APIServer) HandleGetPosts(w http.ResponseWriter, r *http.Request) error {
	resCh := make(chan ThirdPartyResponse)

	go func() {
		fmt.Println("some text to console")

		val, err := thirdPartyApiCall()
		resCh <- ThirdPartyResponse{
			val: val,
			err: err,
		}
	}()

	res := <-resCh // wait for response from third pary
	if res.err != nil {
		fmt.Println("third party failed request")
		return res.err
	}

	posts, err := s.store.GetPosts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, posts)
}

func (s *APIServer) HandleGetPostByID(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return InvalidRequestError("Must provide valid integer for post ID")
	}

	post, err := s.store.GetPostByID(id)
	if err != nil {
		fmt.Println(err.Error())
		return WriteJSON(w, http.StatusNotFound, nil)
	}

	return WriteJSON(w, http.StatusOK, post)
}

func (s *APIServer) HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	newPost := &CreatePostRequest{}
	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return err
	}

	if len(newPost.Body) > 200 {
		return InvalidRequestError("Body needs to be shorter than 200 characters")
	}

	post := NewPost(newPost.Body)
	if err := s.store.CreatePost(post); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, post)
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /", APIHandler(s.HandleGetPosts))
	router.HandleFunc("GET /posts/{id}", APIHandler(s.HandleGetPostByID))
	router.HandleFunc("POST /posts", APIHandler(s.HandleCreatePost))

	m := UseMiddleware(
		LoggingMiddleware,
	)

	server := http.Server{
		Addr:    s.listenAddr,
		Handler: m(router),
	}

	return server.ListenAndServe()
}
