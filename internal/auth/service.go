package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err // DB Error
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	token := hex.EncodeToString(b)

	now := time.Now().UTC()
	expiresAt := now.Add(24 * time.Hour)

	sessions := &Session{
		Token:     token,
		UserId:    user.Id,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	err = s.repo.CreateSession(ctx, sessions)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, ExpiresAt: expiresAt, User: user}, nil
}
