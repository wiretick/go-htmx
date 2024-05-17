package data

type Post struct {
	Id   int
	Body string
}

func NewPost(body string) Post {
	return Post{Id: 1, Body: body}
}
