package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/services"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/httpcommon"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	service *services.Service
}

func New(srv services.Service) *Handler {
	return &Handler{
		service: &srv,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlistID, err := uuid.Parse(r.PathValue("wishlistId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	var input CreateGiftInput
	if err := httpcommon.DecodeJSON(r, &input); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Description) == "" || strings.TrimSpace(input.Link) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	id, err := h.service.Create(ctx, services.CreateGiftInput{
		UserID:      userID,
		WishlistID:  wishlistID,
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		Priority:    input.Priority,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidGift):
			httpcommon.WriteInvalidRequest(w)
		case errors.Is(err, domain.ErrGiftAlreadyExists):
			httpcommon.WriteConflict(w, "gift already exists")
		case errors.Is(err, domain.ErrGiftNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "wishlist not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, CreateGiftOutput{ID: id})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlistID, err := uuid.Parse(r.PathValue("wishlistId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	gifts, err := h.service.ListByWishlistID(ctx, wishlistID, userID)
	if err != nil {
		httpcommon.WriteInternalError(w)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, toGiftsResponse(gifts))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlistID, err := uuid.Parse(r.PathValue("wishlistId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	giftID, err := uuid.Parse(r.PathValue("giftId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	gift, err := h.service.GetByID(ctx, giftID, wishlistID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrGiftNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "gift not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, toGiftResponse(gift))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlistID, err := uuid.Parse(r.PathValue("wishlistId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	giftID, err := uuid.Parse(r.PathValue("giftId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	var input UpdateGiftInput
	if err := httpcommon.DecodeJSON(r, &input); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Description) == "" || strings.TrimSpace(input.Link) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	err = h.service.Update(ctx, services.UpdateGiftInput{
		ID:          giftID,
		UserID:      userID,
		WishlistID:  wishlistID,
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		Priority:    input.Priority,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidGift):
			httpcommon.WriteInvalidRequest(w)
		case errors.Is(err, domain.ErrGiftNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "gift not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlistID, err := uuid.Parse(r.PathValue("wishlistId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	giftID, err := uuid.Parse(r.PathValue("giftId"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	if err := h.service.Delete(ctx, giftID, wishlistID, userID); err != nil {
		switch {
		case errors.Is(err, domain.ErrGiftNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "gift not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, nil)
}
