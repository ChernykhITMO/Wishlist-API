package domain

import "github.com/google/uuid"

type Gift struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	WishlistID  uuid.UUID
	Name        string
	Description string
	Link        string
	Priority    int
}
