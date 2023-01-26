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
    VALUES($1, $2, current_timestamp, current_timestamp + $3::interval)`
	var id int
	err := s.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT * FROM snippets WHERE expires>current_timestamp and id=$1`
	row := s.DB.QueryRow(stmt, id)
	snippet := &models.Snippet{}
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return snippet, nil
}

func (s *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT * FROM snippets
			WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`

	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := make([]*models.Snippet, 0)
	for rows.Next() {
		snippet := &models.Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return snippets, nil
}
