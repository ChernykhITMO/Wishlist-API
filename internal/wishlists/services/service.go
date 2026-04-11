package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
	"github.com/google/uuid"
)

type Service struct {
	wishlistRepo Repository
}

func New(repo Repository) *Service {
	return &Service{
		wishlistRepo: repo,
	}
}

func (s *Service) Create(ctx context.Context, input CreateWishlistInput) (string, error) {
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)
	if name == "" || description == "" {
		return "", domain.ErrInvalidWishlist
	}

	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return "", domain.ErrInvalidWishlist
	}

	wishlist := domain.Wishlist{
		ID:          uuid.New(),
		UserID:      input.UserID,
		Token:       uuid.New(),
		NameEvent:   name,
		Description: description,
		DateEvent:   date.UTC(),
	}

	if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
		if errors.Is(err, domain.ErrWishlistAlreadyExists) {
			return "", domain.ErrWishlistAlreadyExists
		}
		return "", err
	}

	return wishlist.Token.String(), nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	return s.wishlistRepo.ListByUserID(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID) (domain.Wishlist, error) {
	return s.wishlistRepo.GetByID(ctx, id, userID)
}

func (s *Service) GetByToken(ctx context.Context, token uuid.UUID) (domain.Wishlist, error) {
	return s.wishlistRepo.GetByToken(ctx, token)
}

func (s *Service) Update(ctx context.Context, input UpdateWishlistInput) error {
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)
	if name == "" || description == "" {
		return domain.ErrInvalidWishlist
	}

	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return domain.ErrInvalidWishlist
	}

	wishlist := domain.Wishlist{
		ID:          input.ID,
		UserID:      input.UserID,
		NameEvent:   name,
		Description: description,
		DateEvent:   date.UTC(),
	}

	if err := s.wishlistRepo.Update(ctx, wishlist); err != nil {
		if errors.Is(err, domain.ErrWishlistNotFound) {
			return domain.ErrWishlistNotFound
		}
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if err := s.wishlistRepo.Delete(ctx, id, userID); err != nil {
		if errors.Is(err, domain.ErrWishlistNotFound) {
			return domain.ErrWishlistNotFound
		}
		return err
	}

	return nil
}
