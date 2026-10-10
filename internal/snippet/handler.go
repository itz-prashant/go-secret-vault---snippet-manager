package snippet

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/itz-prashant/secret-vault-api/internal/appcontext"
	"github.com/itz-prashant/secret-vault-api/internal/utils/response"
)

type Hanlder struct {
	service *Service
}

func NewHandler(service *Service) *Hanlder {
	return &Hanlder{
		service: service,
	}
}

func (h *Hanlder) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateSnippetRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadGateway, "invalid request payload")
		return
	}

	snippet, err := h.service.Create(r.Context(), userID, req)

	if err != nil {
		if errors.Is(err, ErrTitleRequired) || errors.Is(err, ErrContentRequired) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusCreated, snippet)
}

func (h *Hanlder) HandleGet(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid snippet id")
		return
	}

	snippet, err := h.service.GetById(r.Context(), id, userID)

	if err != nil {
		if errors.Is(err, ErrSnippetNotFound) {
			response.WriteError(w, http.StatusNotFound, "snippet not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJson(w, http.StatusOK, snippet)
}

func (h *Hanlder) HandleList(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	snippets, err := h.service.List(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJson(w, http.StatusOK, snippets)
}

func (h *Hanlder) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())

	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid snippet id")
		return
	}

	var req UpdateSnippetRequest

	if err := response.ReadJson(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	snippet, err := h.service.Update(r.Context(), id, userID, req)

	if err != nil {
		if errors.Is(err, ErrSnippetNotFound) {
			response.WriteError(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		if errors.Is(err, ErrTitleRequired) || errors.Is(err, ErrContentRequired) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJson(w, http.StatusOK, snippet)
}

func (h *Hanlder) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := appcontext.GetUserId(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid snippet id")
		return
	}

	err = h.service.Delete(r.Context(), id, userID)

	if err != nil {
		if errors.Is(err, ErrSnippetNotFound) {
			response.WriteError(w, http.StatusNotFound, "snippet not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
