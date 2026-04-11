package handler

import "github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"

func toGiftResponse(gift domain.Gift) Gift {
	return Gift{
		ID:          gift.ID,
		WishlistID:  gift.WishlistID,
		Name:        gift.Name,
		Description: gift.Description,
		Link:        gift.Link,
		Priority:    gift.Priority,
	}
}

func toGiftsResponse(gifts []domain.Gift) GiftsResponse {
	res := GiftsResponse{
		Gifts: make([]Gift, 0, len(gifts)),
	}
	for _, gift := range gifts {
		res.Gifts = append(res.Gifts, toGiftResponse(gift))
	}
	return res
}
