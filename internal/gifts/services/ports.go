package services

import (
	"context"
	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/google/uuid"
)

type OwnerRepository interface {
	Create(ctx context.Context, gift domain.Gift) error
	GetByID(ctx context.Context, giftID, wishlistID, userID uuid.UUID) (domain.Gift, error)
	ListByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]domain.Gift, error)
	Update(ctx context.Context, gift domain.Gift) error
	Delete(ctx context.Context, giftID, wishlistID, userID uuid.UUID) error
}

type PublicRepository interface {
	ListByToken(ctx context.Context, token uuid.UUID) ([]domain.Gift, error)
}
