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

// Create godoc
// @Summary Create gift
// @Tags gifts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistId path string true "Wishlist ID"
// @Param request body CreateGiftInput true "Gift payload"
// @Success 201 {object} CreateGiftOutput
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 404 {object} httpcommon.ErrorPayload
// @Failure 409 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/{wishlistId}/gifts [post]
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

// List godoc
// @Summary List wishlist gifts
// @Tags gifts
// @Produce json
// @Security BearerAuth
// @Param wishlistId path string true "Wishlist ID"
// @Success 200 {object} GiftsResponse
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/{wishlistId}/gifts [get]
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

// Get godoc
// @Summary Get gift by ID
// @Tags gifts
// @Produce json
// @Security BearerAuth
// @Param wishlistId path string true "Wishlist ID"
// @Param giftId path string true "Gift ID"
// @Success 200 {object} Gift
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 404 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/{wishlistId}/gifts/{giftId} [get]
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

// Update godoc
// @Summary Update gift
// @Tags gifts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistId path string true "Wishlist ID"
// @Param giftId path string true "Gift ID"
// @Param request body UpdateGiftInput true "Gift payload"
// @Success 200
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 404 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/{wishlistId}/gifts/{giftId} [put]
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

// Delete godoc
// @Summary Delete gift
// @Tags gifts
// @Produce json
// @Security BearerAuth
// @Param wishlistId path string true "Wishlist ID"
// @Param giftId path string true "Gift ID"
// @Success 200
// @Failure 400 {object} httpcommon.ErrorPayload
// @Failure 401 {object} httpcommon.ErrorPayload
// @Failure 404 {object} httpcommon.ErrorPayload
// @Failure 500 {object} httpcommon.ErrorPayload
// @Router /wishlists/{wishlistId}/gifts/{giftId} [delete]
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
