package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	if strings.TrimSpace(req.UserName) == "" {
		return nil, ErrUsernameRequired
	}

	if strings.TrimSpace(req.Email) == "" {
		return nil, ErrEmailRequired
	}

	if len(req.Password) < 6 {
		return nil, ErrPasswordTooShort
	}

	haspss, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to bcrypt password %w", err)
	}

	user := &User{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: string(haspss),
		CreatedAt:    time.Now().UTC(),
	}

	err = s.repo.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}
