package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/google/uuid"
)

type ownerRepositoryStub struct {
	createFn func(ctx context.Context, gift domain.Gift) error
}

func (s ownerRepositoryStub) Create(ctx context.Context, gift domain.Gift) error {
	if s.createFn != nil {
		return s.createFn(ctx, gift)
	}
	return nil
}

func (s ownerRepositoryStub) GetByID(ctx context.Context, giftID, wishlistID, userID uuid.UUID) (domain.Gift, error) {
	return domain.Gift{}, nil
}

func (s ownerRepositoryStub) ListByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]domain.Gift, error) {
	return nil, nil
}

func (s ownerRepositoryStub) Update(ctx context.Context, gift domain.Gift) error {
	return nil
}

func (s ownerRepositoryStub) Delete(ctx context.Context, giftID, wishlistID, userID uuid.UUID) error {
	return nil
}

type publicRepositoryStub struct{}

func (publicRepositoryStub) ListByToken(ctx context.Context, token uuid.UUID) ([]domain.Gift, error) {
	return nil, nil
}

func TestCreateTrimsGiftFieldsBeforePersisting(t *testing.T) {
	var createdGift domain.Gift

	svc := New(
		ownerRepositoryStub{
			createFn: func(ctx context.Context, gift domain.Gift) error {
				createdGift = gift
				return nil
			},
		},
		publicRepositoryStub{},
	)

	id, err := svc.Create(context.Background(), CreateGiftInput{
		UserID:      uuid.New(),
		WishlistID:  uuid.New(),
		Name:        "  Massage Gun  ",
		Description: "  Relaxing gadget  ",
		Link:        "  https://example.com/item  ",
		Priority:    3,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if id == uuid.Nil {
		t.Fatal("expected non-empty gift ID")
	}
	if createdGift.Name != "Massage Gun" {
		t.Fatalf("expected trimmed name, got %q", createdGift.Name)
	}
	if createdGift.Description != "Relaxing gadget" {
		t.Fatalf("expected trimmed description, got %q", createdGift.Description)
	}
	if createdGift.Link != "https://example.com/item" {
		t.Fatalf("expected trimmed link, got %q", createdGift.Link)
	}
}

func TestCreateRejectsPriorityOutOfRange(t *testing.T) {
	svc := New(ownerRepositoryStub{}, publicRepositoryStub{})

	_, err := svc.Create(context.Background(), CreateGiftInput{
		UserID:      uuid.New(),
		WishlistID:  uuid.New(),
		Name:        "Massage Gun",
		Description: "Relaxing gadget",
		Link:        "https://example.com/item",
		Priority:    6,
	})
	if !errors.Is(err, domain.ErrInvalidGift) {
		t.Fatalf("expected ErrInvalidGift, got %v", err)
	}
}

func TestCreateMapsWishlistNotFound(t *testing.T) {
	svc := New(
		ownerRepositoryStub{
			createFn: func(ctx context.Context, gift domain.Gift) error {
				return domain.ErrGiftNotFound
			},
		},
		publicRepositoryStub{},
	)

	_, err := svc.Create(context.Background(), CreateGiftInput{
		UserID:      uuid.New(),
		WishlistID:  uuid.New(),
		Name:        "Massage Gun",
		Description: "Relaxing gadget",
		Link:        "https://example.com/item",
		Priority:    3,
	})
	if !errors.Is(err, domain.ErrGiftNotFound) {
		t.Fatalf("expected ErrGiftNotFound, got %v", err)
	}
}
