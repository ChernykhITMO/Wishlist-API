package services

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, id uuid.UUID, email, password string) error
	GetByEmail(ctx context.Context, email string) (uuid.UUID, string, error)
}

type TokenManager interface {
	Issue(userID uuid.UUID) (string, error)
}

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
