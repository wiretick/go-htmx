package data

type Post struct {
	Id   int
	Body string
}

func NewPost(body string) Post {
	return Post{Body: body}
}
