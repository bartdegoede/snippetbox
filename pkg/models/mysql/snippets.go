package mysql

import (
	"database/sql"

	"bartdegoe.de/snippetbox/pkg/models"
)

// define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// this wil insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// this will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// this will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}