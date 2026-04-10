package handler

import (
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/httpcommon"
	"github.com/ChernykhITMO/Wishlist-API/internal/services"
	"net/http"
	"strings"
)

type Handler struct {
	service *services.Service
}

func New(srv services.Service) *Handler {
	return &Handler{
		service: &srv,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	ctx := r.Context()

	if err := httpcommon.DecodeJSON(r, &req); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	token, err := h.service.Register(ctx, services.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyExists):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.CodeInvalidRequest, "email already exists")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, registerResponse{Token: token})
}
