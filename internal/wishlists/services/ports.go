package services

import (
	"context"
	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, wishlist domain.Wishlist) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (domain.Wishlist, error)
	GetByToken(ctx context.Context, token uuid.UUID) (domain.Wishlist, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error)
	Update(ctx context.Context, wishlist domain.Wishlist) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}
