package services

import (
	"context"
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/booking/domain"
	"github.com/google/uuid"
)

type Service struct {
	bookingRepo Repository
}

func New(repo Repository) *Service {
	return &Service{
		bookingRepo: repo,
	}
}

func (s *Service) CreateBooking(ctx context.Context, token, giftID uuid.UUID) error {
	bookingID := uuid.New()
	if err := s.bookingRepo.Create(ctx, token, giftID, bookingID); err != nil {
		if errors.Is(err, domain.ErrBookingAlreadyExists) {
			return domain.ErrBookingAlreadyExists
		}
		return err
	}

	return nil
}
