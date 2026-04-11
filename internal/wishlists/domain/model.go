package domain

import (
	"github.com/google/uuid"
	"time"
)

type Wishlist struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Token       uuid.UUID
	NameEvent   string
	Description string
	DateEvent   time.Time
}
