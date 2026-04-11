package domain

import "errors"

var (
	ErrGiftAlreadyExists = errors.New("gift already exists")
	ErrGiftNotFound      = errors.New("gift not found")
	ErrInvalidGift       = errors.New("invalid gift")
)
