package domain

import "errors"

var (
	ErrBookingAlreadyExists = errors.New("booking already exists")
	ErrGiftNotFound         = errors.New("gift not found")
)
