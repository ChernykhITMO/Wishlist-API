package services

import (
	"github.com/google/uuid"
)

type CreateWishlistInput struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
}

type UpdateWishlistInput struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
}
