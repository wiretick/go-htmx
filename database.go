package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetPosts() ([]*Post, error)
	GetPostByID(id int) (*Post, error)
	CreatePost(*Post) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// Please don't hack my local postgresdatase :*)
	connStr := "postgres://postgres:gohtmx@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	if err := s.createPostTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) createPostTable() error {
	log.Println("Creating posts table if it does not already exist")

	query := `CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		body TEXT NOT NULL,
		created_at TIMESTAMP default timezone('utc', now())
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetPosts() ([]*Post, error) {
	query := `SELECT * FROM posts`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []*Post{}
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.Id, &post.Body, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresStore) GetPostByID(id int) (*Post, error) {
	post := &Post{}

	query := `SELECT * FROM posts WHERE id=$1`
	if err := s.db.QueryRow(query, id).Scan(&post.Id, &post.Body, &post.CreatedAt); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostgresStore) CreatePost(post *Post) error {
	var id int

	query := `INSERT INTO posts (body) VALUES ($1) RETURNING id`
	if err := s.db.QueryRow(query, post.Body).Scan(&id); err != nil {
		return err
	}

	log.Println("new post id: ", id)
	return nil
}
