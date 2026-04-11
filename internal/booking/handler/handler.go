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

// CreateBooking godoc
// @Summary Create public gift booking
// @Tags public
// @Accept json
// @Produce json
// @Param token path string true "Public token"
// @Param request body CreateBookingRequest true "Booking payload"
// @Success 201
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 409 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/public/{token}/bookings [post]
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
