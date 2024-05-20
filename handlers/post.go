package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wiretick/go-htmx/core"
	"github.com/wiretick/go-htmx/data"
	"github.com/wiretick/go-htmx/view/post"
)

// Create gaming sessions
var posts []data.Post = []data.Post{
	{Body: "whats long text"},
	{Body: "body another longer text"},
	{Body: "Need to move on to the next topic now"},
}

type ThirdPartyResponse struct {
	val string
	err error
}

func thirdPartyApiCall() (string, error) {
	time.Sleep(time.Millisecond * 500)
	return "data retrieved", nil
}

func HandleGetPosts(w http.ResponseWriter, r *http.Request) error {
	if err := post.Index(posts).Render(r.Context(), w); err != nil {
		return err
	}
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return core.InvalidRequestError("Must provide valid integer for post ID")
	}

	if id < 0 || id >= len(posts) {
		return core.NotFoundError("Could not find a post with the given ID")
	}

	if err := core.WriteJSON(w, http.StatusOK, posts[id]); err != nil {
		return err
	}

	return nil
}

func HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	title := r.FormValue("title")
	body := r.FormValue("body")

	if len(title) > 20 {
		return core.InvalidRequestError("Title needs to be shorter than 20 characters")
	}

	p := data.NewPost(body)
	posts = append(posts, p)

	w.WriteHeader(http.StatusSeeOther)
	return nil
}
