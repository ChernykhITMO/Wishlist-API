package services

import (
	"context"
	"testing"

	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
	"github.com/google/uuid"
)

type repositoryStub struct {
	createFn func(ctx context.Context, wishlist domain.Wishlist) error
}

func (s repositoryStub) Create(ctx context.Context, wishlist domain.Wishlist) error {
	if s.createFn != nil {
		return s.createFn(ctx, wishlist)
	}
	return nil
}

func (s repositoryStub) GetByID(ctx context.Context, id, userID uuid.UUID) (domain.Wishlist, error) {
	return domain.Wishlist{}, nil
}

func (s repositoryStub) GetByToken(ctx context.Context, token uuid.UUID) (domain.Wishlist, error) {
	return domain.Wishlist{}, nil
}

func (s repositoryStub) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	return nil, nil
}

func (s repositoryStub) Update(ctx context.Context, wishlist domain.Wishlist) error {
	return nil
}

func (s repositoryStub) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return nil
}

func TestCreateReturnsWishlistIDAndToken(t *testing.T) {
	var created domain.Wishlist

	svc := New(repositoryStub{
		createFn: func(ctx context.Context, wishlist domain.Wishlist) error {
			created = wishlist
			return nil
		},
	})

	result, err := svc.Create(context.Background(), CreateWishlistInput{
		UserID:      uuid.New(),
		Name:        "Birthday",
		Description: "Gifts for birthday",
		Date:        "2026-12-31T21:00:00+03:00",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if result.ID == uuid.Nil {
		t.Fatal("expected non-empty wishlist ID")
	}
	if result.Token == uuid.Nil {
		t.Fatal("expected non-empty wishlist token")
	}
	if result.ID != created.ID {
		t.Fatalf("expected result ID %s, got %s", created.ID, result.ID)
	}
	if result.Token != created.Token {
		t.Fatalf("expected result token %s, got %s", created.Token, result.Token)
	}
}
