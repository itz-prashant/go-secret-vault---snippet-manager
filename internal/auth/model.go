package auth

import (
	"context"
	"errors"
	"time"
)

var ErrUsernameRequired = errors.New("username is required")
var ErrUserNotFound = errors.New("user not found")
var ErrEmailRequired = errors.New("email is required")
var ErrPasswordTooShort = errors.New("password must be atleast 6 character")
var ErrUserAlreadyExists = errors.New("username and email already exist")
var ErrInvalidCredentials = errors.New("invalid email or password")

type User struct {
	Id int64 `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Token string `json:"token"`
	UserId int64 `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User *User `json:"user"`
}

type Repository interface {
	CreateUser (ctx context.Context, user *User) error
	GetUserByEmail (ctx context.Context, email string) (*User, error)
	GetUserByUserName (ctx context.Context, username string) (*User, error)
	CreateSession (ctx context.Context, session *Session)  error
}