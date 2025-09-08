// Filename: internal/data/comments.go
package data

import (
	"time"
	"github.com/amilcar-vasquez/qod/internal/validator"
)

// make our JSON keys be displayed in all lowercase
// "-" means don't show this field
type Comment struct {
    ID int64                     `json:"id"`                   
    Content  string              `json:"content"`     
    Author  string               `json:"author"`
    CreatedAt  time.Time         `json:"-"`     
    Version int32                `json:"version"`      
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
