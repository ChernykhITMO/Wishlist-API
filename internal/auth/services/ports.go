package postgres

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, id uuid.UUID, email, password string) error
	GetUserByEmail(ctx context.Context, email string) (uuid.UUID, string, error)
}
