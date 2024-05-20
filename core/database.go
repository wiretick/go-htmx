package core

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/wiretick/go-htmx/data"
)

type Storage interface {
	CreatePost(*data.Post) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// Please don't hack my local postgres database :*)
	connStr := "postgres://postgres:gohtmx@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &PostgresStore{}, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createPostTable()
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

func (s *PostgresStore) CreatePost(post *data.Post) error {
	var id int

	query := `INSERT INTO posts (body) VALUES ($1) RETURNING id`
	if err := s.db.QueryRow(query, post.Body).Scan(&id); err != nil {
		return err
	}

	log.Println("new post id: ", id)
	return nil
}
