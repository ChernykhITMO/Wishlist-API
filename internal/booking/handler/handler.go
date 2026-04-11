package handler

import (
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/booking/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/booking/services"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/httpcommon"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	service *services.Service
}

func New(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, err := uuid.Parse(r.PathValue("token"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	var req CreateBookingRequest
	if err := httpcommon.DecodeJSON(r, &req); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	if err := h.service.CreateBooking(ctx, token, req.GiftID); err != nil {
		switch {
		case errors.Is(err, domain.ErrBookingAlreadyExists):
			httpcommon.WriteError(w, http.StatusConflict, httpcommon.CodeConflict, "booking already exists")
		case errors.Is(err, domain.ErrGiftNotFound):
			httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.CodeGiftNotFound, "gift not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, nil)
}
