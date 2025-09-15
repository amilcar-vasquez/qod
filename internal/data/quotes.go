// Filename: internal/data/quotes.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/amilcar-vasquez/qod/internal/validator"
)

// A QuoteModel expects a connection pool
type QuoteModel struct {
	DB *sql.DB
}

// make our JSON keys be displayed in all lowercase
// "-" means don't show this field
type Quote struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

// Insert a new row in the quotes table
// Expects a pointer to the actual quote
func (q QuoteModel) Insert(quote *Quote) error {
	query := `
	INSERT INTO qod (content, author)
	VALUES ($1, $2)
	RETURNING id, created_at, version
	`
	args := []any{quote.Content, quote.Author}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return q.DB.QueryRowContext(ctx, query, args...).Scan(
		&quote.ID,
		&quote.CreatedAt,
		&quote.Version)
}

// Get a specific quote based on its ID
func (q QuoteModel) Get(id int64) (*Quote, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, content, author, created_at, version
	FROM qod
	WHERE id = $1`

	var quote Quote
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := q.DB.QueryRowContext(ctx, query, id).Scan(
		&quote.ID,
		&quote.Content,
		&quote.Author,
		&quote.CreatedAt,
		&quote.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &quote, nil
}

// update a specific quote based on its ID
func (q QuoteModel) Update(quote *Quote) error {
	query := `
	UPDATE qod
	SET content = $1, author = $2, version = version + 1
	WHERE id = $3
	RETURNING version`
	args := []any{
		quote.Content,
		quote.Author,
		quote.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)
}

// delete a specific quote based on its ID
func (q QuoteModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM qod
	WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := q.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Get all the quotes
func (q QuoteModel) GetAll(content string, author string, filters Filters) ([]*Quote, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, content, author, created_at, version
	FROM qod
	WHERE (to_tsvector('simple', content) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', author) @@ plainto_tsquery('simple', $2) OR $2 = '')
	 ORDER BY %s %s, id ASC 
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := q.DB.QueryContext(ctx, query, content, author, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		err := rows.Scan(&totalRecords,
			&quote.ID,
			&quote.Content,
			&quote.Author,
			&quote.CreatedAt,
			&quote.Version)
		if err != nil {
			return nil, Metadata{}, err
		}
		quotes = append(quotes, &quote)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return quotes, metadata, nil
}

// Create a function that performs the validation checks
func ValidateQuote(v *validator.Validator, quote *Quote) {
	v.Check(quote.Content != "", "content", "must be provided")
	v.Check(quote.Author != "", "author", "must be provided")
	v.Check(len(quote.Content) <= 100, "content", "must not be more than 100 bytes long")
	v.Check(len(quote.Author) <= 25, "author", "must not be more than 25 bytes long")
}
