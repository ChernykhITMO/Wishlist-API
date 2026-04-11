package domain

import "errors"

var (
	ErrWishlistAlreadyExists = errors.New("wishlist already exists")
	ErrWishlistNotFound      = errors.New("wishlist not found")
	ErrInvalidWishlist       = errors.New("invalid wishlist")
)
