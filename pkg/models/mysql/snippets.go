package mysql

import (
	"database/sql"
	"errors"

	"bartdegoe.de/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &models.Snippet{}

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows result set.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows result set is
	// always properly closed before the Latest() method returns. This deferr
	// statement should come *after* the check for an error from the Query()
	// method. Otherwise, if Query() returns an error, we'll generate a panic
	// when we try to close a nil result set.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippet objects
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the result set. This
	// prepares the first (and then subsequen) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes, then
	// the result set automatically closes itself and frees up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed models.Snippet object
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append each hydrated snippet to the result slice
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed over
	// the whole result set.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything's okay, return the result set
	return snippets, nil
}
