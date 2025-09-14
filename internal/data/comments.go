// Filename: internal/data/comments.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/amilcar-vasquez/qod/internal/validator"
	"time"
)

// A CommentModel expects a connection pool
type CommentModel struct {
	DB *sql.DB
}

// make our JSON keys be displayed in all lowercase
// "-" means don't show this field
type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

// Insert a new row in the comments table
// Expects a pointer to the actual comment
func (c CommentModel) Insert(comment *Comment) error {
	// the SQL query to be executed against the database table
	query := `
        INSERT INTO qod (content, author)
        VALUES ($1, $2)
        RETURNING id, created_at, version
        `
	// the actual values to replace $1, and $2
	args := []any{comment.Content, comment.Author}
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute the query against the comments database table. We ask for the the
	// id, created_at, and version to be sent back to us which we will use
	// to update the Comment struct later on
	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Version)
}

// Get a specific comment based on its ID
func (c CommentModel) Get(id int64) (*Comment, error) {
	// if the id is less than 1, then it's an invalid ID
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, content, author, created_at, version
		FROM qod
		WHERE id = $1`

	//declare a variable to hold the returned data
	var comment Comment
	// Set a 3-second context/timer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.Author,
		&comment.CreatedAt,
		&comment.Version)
	//check for error that may have occurred when executing the query
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	//cancel the context before returning
	return &comment, nil
}

// update a specific comment based on its ID
func (c CommentModel) Update(comment *Comment) error {
	// the SQL query to be executed against the database table
	query := `
		UPDATE qod
		SET content = $1, author = $2, version = version + 1
		WHERE id = $3
		RETURNING version`
	// the actual values to replace $1, $2, and $3
	args := []any{
		comment.Content,
		comment.Author,
		comment.ID,
	}
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute the query against the comments database table. We ask for the the
	// version to be sent back to us which we will use to update the Comment struct later on

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
}

// delete a specific comment based on its ID
func (c CommentModel) Delete(id int64) error {
	// if the id is less than 1, then it's an invalid ID
	if id < 1 {
		return ErrRecordNotFound
	}
	// the SQL query to be executed against the database table
	query := `
		DELETE FROM qod
		WHERE id = $1`
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ExecContext does not return any rows unlike QueryRowContext.
	// It only returns  information about the the query execution
	// such as how many rows were affected
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Get all the comments
func (c CommentModel) GetAll() ([]*Comment, error) {
	// the SQL query to be executed against the database table
	query := `
		SELECT id, content, author, created_at, version
		FROM qod
		ORDER BY id`
	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute the query against the comments database table
	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// close the result set before the GetAll() method exits
	defer rows.Close()
	// initialize an empty slice to hold the comment data
	var comments []*Comment
	// iterate through the rows in the result set
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
			&comment.Version)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	// check for errors after we are done iterating through the rows
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

// Create a function that performs the validation checks
func ValidateComment(v *validator.Validator, comment *Comment) {
	// check if the Content field is empty
	v.Check(comment.Content != "", "content", "must be provided")
	// check if the Author field is empty
	v.Check(comment.Author != "", "author", "must be provided")
	// check if the Content field is empty
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
}
