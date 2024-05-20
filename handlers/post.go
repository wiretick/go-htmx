package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wiretick/go-htmx/core"
	"github.com/wiretick/go-htmx/data"
)

type ThirdPartyResponse struct {
	val string
	err error
}

func thirdPartyApiCall() (string, error) {
	time.Sleep(time.Millisecond * 500)
	return "data retrieved", nil
}

func HandleGetPosts(w http.ResponseWriter, r *http.Request) error {
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

	fmt.Println(res.val)

	return nil
}

func HandleGetPostByID(w http.ResponseWriter, r *http.Request) error {
	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return core.InvalidRequestError("Must provide valid integer for post ID")
	}

	//if id < 0 || id >= len(posts) {
	//	return core.NotFoundError("Could not find a post with the given ID")
	//}

	//if err := core.WriteJSON(w, http.StatusOK, posts[id]); err != nil {
	//	return err
	//}

	return nil
}

func (s *APIServer) HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	newPost := &data.CreatePostRequest{}
	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return err
	}

	if len(newPost.Body) > 200 {
		return core.InvalidRequestError("Body needs to be shorter than 200 characters")
	}

	post := data.NewPost(newPost.Body)
	if err := s.store.CreatePost(post); err != nil {
		return err
	}

	return core.WriteJSON(w, http.StatusOK, post)
}

type APIServer struct {
	listenAddr string
	store      core.Storage
}

func NewAPIServer(listenAddr string, store core.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", core.APIHandler(HandleGetPosts))
	router.HandleFunc("GET /posts/{id}", core.APIHandler(HandleGetPostByID))
	router.HandleFunc("POST /posts", core.APIHandler(s.HandleCreatePost))

	m := core.UseMiddleware(
		core.LoggingMiddleware,
	)

	server := http.Server{
		Addr:    s.listenAddr,
		Handler: m(router),
	}

	log.Fatal(server.ListenAndServe())
}
