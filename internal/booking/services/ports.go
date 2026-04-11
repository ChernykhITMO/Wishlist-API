package internal

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, token, giftID, bookingID uuid.UUID) error
}
