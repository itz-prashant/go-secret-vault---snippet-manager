package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/itz-prashant/secret-vault-api/internal/appcontext"
	"github.com/itz-prashant/secret-vault-api/internal/auth"
	"github.com/itz-prashant/secret-vault-api/internal/utils/response"
)


type AuthMiddleware struct {
	authRepo auth.Repository
}

func NewAuthMiddleware (authRepo auth.Repository) *AuthMiddleware{
	return &AuthMiddleware{
		authRepo: authRepo,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")

		if authHeader == ""{
			response.WriteError(w, http.StatusUnauthorized, "authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteError(w, http.StatusUnauthorized, "invalid authorization format: expected 'Bearer <token>'")
			return 
		}
		token := parts[1]

		session, err:= m.authRepo.GetSession(r.Context(), token)

		if err != nil {
			if errors.Is(err, auth.ErrSessionNotFound){
				response.WriteError(w, http.StatusUnauthorized, "invaid or expired token")
				return
			}
			response.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if session.ExpiresAt.Before(time.Now().UTC()){
			response.WriteError(w, http.StatusUnauthorized, "session has expired, please login")
			return
		}

		ctx := appcontext.SetUserId(r.Context(), session.UserId)

		next(w, r.WithContext(ctx))
	}
}