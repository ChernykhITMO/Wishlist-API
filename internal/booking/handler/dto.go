package handler

import "github.com/google/uuid"

type CreateBookingRequest struct {
	GiftID uuid.UUID `json:"giftId"`
}
