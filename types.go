package main

import "time"

type Post struct {
	Id        int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreatePostRequest struct {
	Body string `json:"body"`
}

func NewPost(body string) *Post {
	return &Post{
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
}
