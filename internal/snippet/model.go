package snippet

import (
	"context"
	"time"
)

type Snippet struct {
	ID        int16     `json:"id"`
	UserId    int16     `json:"user_id"`
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
}