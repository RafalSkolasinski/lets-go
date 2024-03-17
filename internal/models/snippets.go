package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for and individual
// Notice how the fields of the struct corresponds to the fields
// in our MySQL snippets table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps an sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute.
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//Use the Exec() method on the embedded connection pool to execute the statement.
	// The first parameter is the SQL statement, followed by the title, content and expiry
	// values for the placeholder parameters. This method returns an sql.Result type, which
	// which contains some basic information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to g et the ID of our newly inserted record
	// in the snippets table. NOTE: PostgreSQL does not support LastInsertId() method!!!
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned is type int64, convert to int type
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return &Snippet{}, nil
	// return &Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
