package services

import "github.com/google/uuid"

type CreateGiftInput struct {
	UserID      uuid.UUID
	WishlistID  uuid.UUID
	Name        string
	Description string
	Link        string
	Priority    int
}

type UpdateGiftInput struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	WishlistID  uuid.UUID
	Name        string
	Description string
	Link        string
	Priority    int
}
