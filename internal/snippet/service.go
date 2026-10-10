package snippet

import (
	"context"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateSnippetRequest) (*Snippet, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrTitleRequired
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, ErrContentRequired
	}

	isPrivate := true

	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}
	now := time.Now().UTC()

	snippet := &Snippet{
		UserId: userID,
		Title:     req.Title,
		Content:   req.Content,
		IsPrivate: isPrivate,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.Create(ctx, snippet)

	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *Service) GetById(ctx context.Context, id int64, userID int64) (*Snippet, error) {
	snippet, err := s.repo.GetById(ctx, id, userID)

	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]Snippet, error) {
	snippets, err := s.repo.ListByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return snippets, nil
}

func (s *Service) Update(ctx context.Context,id int64, userID int64, req UpdateSnippetRequest) (*Snippet, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, ErrTitleRequired
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, ErrTitleRequired
	}

	snippet, err := s.repo.GetById(ctx, id, userID)

	if err != nil {
		return nil, err
	}

	snippet = &Snippet{
		Title: req.Title,
		Content: req.Content,
		UpdatedAt: time.Now().UTC(),
	}

	err = s.repo.Update(ctx, snippet)

	if err != nil {
		return snippet, err
	}

	return snippet, nil
}

func (s *Service) Delete(ctx context.Context, id int64, userID int64) error {
	err := s.repo.Delete(ctx, id, userID)

	if err != nil {
		return err
	}

	return  nil
}