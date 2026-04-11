package handler

import (
	"context"
	"errors"
	giftdomain "github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/httpcommon"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/middleware"
	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/services"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Handler struct {
	service     *services.Service
	publicGifts publicGiftsLister
}

type publicGiftsLister interface {
	ListByToken(ctx context.Context, token uuid.UUID) ([]giftdomain.Gift, error)
}

func New(srv services.Service, publicGifts publicGiftsLister) *Handler {
	return &Handler{
		service:     &srv,
		publicGifts: publicGifts,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	var input CreateWishlistInput

	if err := httpcommon.DecodeJSON(r, &input); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Description) == "" || strings.TrimSpace(input.Date) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	token, err := h.service.Create(ctx, services.CreateWishlistInput{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
	})

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidWishlist):
			httpcommon.WriteInvalidRequest(w)
		case errors.Is(err, domain.ErrWishlistAlreadyExists):
			httpcommon.WriteConflict(w, "wishlist already exists")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, CreateWishlistOutput{Token: token})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	wishlists, err := h.service.ListByUserID(ctx, userID)
	if err != nil {
		httpcommon.WriteInternalError(w)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, toWishlistsResponse(wishlists))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	wishlist, err := h.service.GetByID(ctx, id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWishlistNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "wishlist not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, toWishlistResponse(wishlist))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		httpcommon.WriteUnauthorized(w, "unauthorized")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	var input UpdateWishlistInput
	if err := httpcommon.DecodeJSON(r, &input); err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Description) == "" || strings.TrimSpace(input.Date) == "" {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	err = h.service.Update(ctx, services.UpdateWishlistInput{
		ID:          id,
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidWishlist):
			httpcommon.WriteInvalidRequest(w)
		case errors.Is(err, domain.ErrWishlistNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "wishlist not found")
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

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	if err := h.service.Delete(ctx, id, userID); err != nil {
		switch {
		case errors.Is(err, domain.ErrWishlistNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "wishlist not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := uuid.Parse(r.PathValue("token"))
	if err != nil {
		httpcommon.WriteInvalidRequest(w)
		return
	}

	wishlist, err := h.service.GetByToken(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWishlistNotFound):
			httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "wishlist not found")
		default:
			httpcommon.WriteInternalError(w)
		}
		return
	}

	gifts, err := h.publicGifts.ListByToken(ctx, token)
	if err != nil {
		httpcommon.WriteInternalError(w)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, toPublicWishlistResponse(wishlist, gifts))
}
