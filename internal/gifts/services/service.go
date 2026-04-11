package services

import (
	"context"
	"errors"
	"strings"

	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/google/uuid"
)

type Service struct {
	ownerRepo  OwnerRepository
	publicRepo PublicRepository
}

func New(ownerRepo OwnerRepository, publicRepo PublicRepository) *Service {
	return &Service{
		ownerRepo:  ownerRepo,
		publicRepo: publicRepo,
	}
}

func (s *Service) Create(ctx context.Context, input CreateGiftInput) (uuid.UUID, error) {
	gift, err := newGift(input)
	if err != nil {
		return uuid.Nil, err
	}

	gift.ID = uuid.New()

	if err := s.ownerRepo.Create(ctx, gift); err != nil {
		switch {
		case errors.Is(err, domain.ErrGiftAlreadyExists):
			return uuid.Nil, domain.ErrGiftAlreadyExists
		case errors.Is(err, domain.ErrGiftNotFound):
			return uuid.Nil, domain.ErrGiftNotFound
		default:
			return uuid.Nil, err
		}
	}

	return gift.ID, nil
}

func (s *Service) ListByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]domain.Gift, error) {
	return s.ownerRepo.ListByWishlistID(ctx, wishlistID, userID)
}

func (s *Service) GetByID(ctx context.Context, giftID, wishlistID, userID uuid.UUID) (domain.Gift, error) {
	gift, err := s.ownerRepo.GetByID(ctx, giftID, wishlistID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrGiftNotFound) {
			return domain.Gift{}, domain.ErrGiftNotFound
		}
		return domain.Gift{}, err
	}

	return gift, nil
}

func (s *Service) Update(ctx context.Context, input UpdateGiftInput) error {
	gift, err := newGift(CreateGiftInput{
		UserID:      input.UserID,
		WishlistID:  input.WishlistID,
		Name:        input.Name,
		Description: input.Description,
		Link:        input.Link,
		Priority:    input.Priority,
	})
	if err != nil {
		return err
	}

	gift.ID = input.ID

	if err := s.ownerRepo.Update(ctx, gift); err != nil {
		if errors.Is(err, domain.ErrGiftNotFound) {
			return domain.ErrGiftNotFound
		}
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, giftID, wishlistID, userID uuid.UUID) error {
	if err := s.ownerRepo.Delete(ctx, giftID, wishlistID, userID); err != nil {
		if errors.Is(err, domain.ErrGiftNotFound) {
			return domain.ErrGiftNotFound
		}
		return err
	}

	return nil
}

func (s *Service) ListByToken(ctx context.Context, token uuid.UUID) ([]domain.Gift, error) {
	gifts, err := s.publicRepo.ListByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return gifts, nil
}

func newGift(input CreateGiftInput) (domain.Gift, error) {
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)
	link := strings.TrimSpace(input.Link)
	if name == "" || description == "" || link == "" {
		return domain.Gift{}, domain.ErrInvalidGift
	}

	if input.Priority < 1 || input.Priority > 5 {
		return domain.Gift{}, domain.ErrInvalidGift
	}

	return domain.Gift{
		UserID:      input.UserID,
		WishlistID:  input.WishlistID,
		Name:        name,
		Description: description,
		Link:        link,
		Priority:    input.Priority,
	}, nil
}
