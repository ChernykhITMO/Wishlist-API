package handler

import (
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/auth/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/auth/services"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/httpcommon"
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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration payload"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 409 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
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
		case errors.Is(err, domain.ErrInvalidRegistration):
			httpcommon.WriteInvalidRequest(w)
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			httpcommon.WriteConflict(w, "email already exists")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, RegisterResponse{Token: token})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	ctx := r.Context()

	if err := httpcommon.DecodeJSON(r, &req); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	token, err := h.service.Login(ctx, services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			httpcommon.WriteUnauthorized(w, "invalid credentials")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, LoginResponse{Token: token})
}
