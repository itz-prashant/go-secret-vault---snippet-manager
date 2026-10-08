package auth

import (
	"errors"
	"net/http"

	"github.com/itz-prashant/secret-vault-api/internal/appcontext"
	"github.com/itz-prashant/secret-vault-api/internal/utils/response"
)

type Handler struct {
	service *Service
}

func NeWHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HanldeRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
	}

	user, err := h.service.Register(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			response.WriteError(w, http.StatusConflict, err.Error())
			return
		}

		if errors.Is(err, ErrEmailRequired) || errors.Is(err, ErrUsernameRequired) || errors.Is(err, ErrPasswordTooShort) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "Internal server error")
	}

	response.WriteJson(w, http.StatusCreated, user)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	loginResp, err := h.service.Login(r.Context(), req)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusOK, loginResp)
}

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorizes")
		return
	}

	response.WriteJson(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"user_id":       userID,
		"message":       "you are accessing a protected route!",
	})
}
