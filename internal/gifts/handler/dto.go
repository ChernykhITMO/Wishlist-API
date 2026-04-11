package handler

import "github.com/google/uuid"

type CreateGiftInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority"`
}

type CreateGiftOutput struct {
	ID uuid.UUID `json:"id"`
}

type UpdateGiftInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority"`
}

type GiftsResponse struct {
	Gifts []Gift `json:"gifts"`
}

type Gift struct {
	ID          uuid.UUID `json:"id"`
	WishlistID  uuid.UUID `json:"wishlistId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Priority    int       `json:"priority"`
}
