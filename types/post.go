package types

type Post struct {
	Title string
	Body  string
}

type PostPage struct {
	Posts Posts
}

type Posts []Post

func (p Posts) Find(id int) Post {
	// pretty useless
	return p[id]
}
