package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/Wishlist-API/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepositoryStub struct {
	createFn     func(ctx context.Context, id uuid.UUID, email, password string) error
	getByEmailFn func(ctx context.Context, email string) (uuid.UUID, string, error)
}

func (s userRepositoryStub) Create(ctx context.Context, id uuid.UUID, email, password string) error {
	if s.createFn != nil {
		return s.createFn(ctx, id, email, password)
	}
	return nil
}

func (s userRepositoryStub) GetByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	if s.getByEmailFn != nil {
		return s.getByEmailFn(ctx, email)
	}
	return uuid.Nil, "", nil
}

type tokenManagerStub struct {
	issueFn func(userID uuid.UUID) (string, error)
}

func (s tokenManagerStub) Issue(userID uuid.UUID) (string, error) {
	if s.issueFn != nil {
		return s.issueFn(userID)
	}
	return "issued-token", nil
}

type passwordManagerStub struct {
	hashFn    func(password string) (string, error)
	compareFn func(hash, password string) error
}

func (s passwordManagerStub) Hash(password string) (string, error) {
	if s.hashFn != nil {
		return s.hashFn(password)
	}
	return "hashed-password", nil
}

func (s passwordManagerStub) Compare(hash, password string) error {
	if s.compareFn != nil {
		return s.compareFn(hash, password)
	}
	return nil
}

func TestRegisterNormalizesEmailBeforePersisting(t *testing.T) {
	var createdEmail string

	svc := New(
		userRepositoryStub{
			createFn: func(ctx context.Context, id uuid.UUID, email, password string) error {
				createdEmail = email
				return nil
			},
		},
		tokenManagerStub{},
		passwordManagerStub{},
	)

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "  USER@Example.COM  ",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if createdEmail != "user@example.com" {
		t.Fatalf("expected normalized email %q, got %q", "user@example.com", createdEmail)
	}
}

func TestRegisterRejectsBlankPassword(t *testing.T) {
	svc := New(userRepositoryStub{}, tokenManagerStub{}, passwordManagerStub{})

	_, err := svc.Register(context.Background(), RegisterInput{
		Email:    "user@example.com",
		Password: "   ",
	})
	if !errors.Is(err, domain.ErrInvalidRegistration) {
		t.Fatalf("expected ErrInvalidRegistration, got %v", err)
	}
}

func TestLoginReturnsInvalidCredentialsWhenUserNotFound(t *testing.T) {
	svc := New(
		userRepositoryStub{
			getByEmailFn: func(ctx context.Context, email string) (uuid.UUID, string, error) {
				return uuid.Nil, "", pgx.ErrNoRows
			},
		},
		tokenManagerStub{},
		passwordManagerStub{},
	)

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
