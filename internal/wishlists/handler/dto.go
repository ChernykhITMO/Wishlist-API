package handler

import (
	giftdomain "github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/google/uuid"
	"time"
)

type CreateWishlistInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type CreateWishlistOutput struct {
	Token string `json:"token"`
}

type UpdateWishlistInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type WishlistsResponse struct {
	Wishlists []Wishlist `json:"wishlists"`
}

type Wishlist struct {
	ID          uuid.UUID `json:"id"`
	NameEvent   string    `json:"nameEvent"`
	Description string    `json:"description"`
	DateEvent   time.Time `json:"dateEvent"`
}

type PublicWishlistResponse struct {
	ID          uuid.UUID       `json:"id"`
	NameEvent   string          `json:"nameEvent"`
	Description string          `json:"description"`
	DateEvent   time.Time       `json:"dateEvent"`
	Gifts       []PublicGiftDTO `json:"gifts"`
}

type PublicGiftDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Priority    int       `json:"priority"`
}

func toPublicGiftDTO(gift giftdomain.Gift) PublicGiftDTO {
	return PublicGiftDTO{
		ID:          gift.ID,
		Name:        gift.Name,
		Description: gift.Description,
		Link:        gift.Link,
		Priority:    gift.Priority,
	}
}
