package handler

import (
	giftdomain "github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
)

func toWishlistResponse(wishlist domain.Wishlist) Wishlist {
	return Wishlist{
		ID:          wishlist.ID,
		NameEvent:   wishlist.NameEvent,
		Description: wishlist.Description,
		DateEvent:   wishlist.DateEvent,
	}
}

func toWishlistsResponse(wishlists []domain.Wishlist) WishlistsResponse {
	res := WishlistsResponse{
		Wishlists: make([]Wishlist, 0, len(wishlists)),
	}
	for _, w := range wishlists {
		res.Wishlists = append(res.Wishlists, toWishlistResponse(w))
	}
	return res
}

func toPublicWishlistResponse(wishlist domain.Wishlist, gifts []giftdomain.Gift) PublicWishlistResponse {
	res := PublicWishlistResponse{
		ID:          wishlist.ID,
		NameEvent:   wishlist.NameEvent,
		Description: wishlist.Description,
		DateEvent:   wishlist.DateEvent,
		Gifts:       make([]PublicGiftDTO, 0, len(gifts)),
	}
	for _, gift := range gifts {
		res.Gifts = append(res.Gifts, toPublicGiftDTO(gift))
	}
	return res
}
