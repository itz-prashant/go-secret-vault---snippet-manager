package snippet

import (
	"context"
	"errors"
	"time"
)

var ErrSnippetNotFound = errors.New("snippet not found")
var ErrTitleRequired = errors.New("title is required")
var ErrContentRequired = errors.New("content is required")

type Snippet struct {
	ID        int64    `json:"id"`
	UserId    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSnippetRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsPrivate *bool   `json:"is_private"`
}

type UpdateSnippetRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsPrivate *bool   `json:"is_private"`
}

type Repository interface {
	Create (ctx context.Context, s *Snippet) error
	GetById (ctx context.Context, id int64, userID int64) (*Snippet, error)
	ListByUserID (ctx context.Context, userID int64) ([]Snippet, error)
	Update (ctx context.Context, snippet *Snippet) error
	Delete (ctx context.Context, id int64, userID int64) error
}