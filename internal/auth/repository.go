package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

func (s *sqliteRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES(?,?,?,?)`

	result, err := s.db.ExecContext(ctx, query, user.UserName, user.Email, user.PasswordHash, user.CreatedAt)

	if err != nil {

		// Unique username/email constraint check
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("Failed to get last insert id %w", err)
	}

	user.Id = id

	return nil
}

func (s *sqliteRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?`

	row := s.db.QueryRowContext(ctx, query, email)

	var user User

	err := row.Scan(&user.Id, &user.UserName, &user.Email, &user.PasswordHash ,&user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (s *sqliteRepository) GetUserByUserName(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?`

	row := s.db.QueryRowContext(ctx, query, username)

	var user User

	err := row.Scan(&user.Id, &user.UserName, &user.Email, &user.PasswordHash , &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)

	}

	return &user, nil
}

func (s *sqliteRepository) CreateSession(ctx context.Context, session *Session) error {
	query := `INSERT INTO sessions (token, user_id, expires_at, created_at) VALUES (?,?,?,?)`

	_, err := s.db.ExecContext(ctx, query, session.Token, session.UserId, session.ExpiresAt, session.CreatedAt)

	if err != nil {
		return  fmt.Errorf("failed to insert session %w", err)
	}

	return nil
}