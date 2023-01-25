package postgresql

import (
	"database/sql"
	"snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES($1, $2, current_timestamp, current_timestamp + $3::interval) RETURNING id`
	var id int
	err := s.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (s *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
