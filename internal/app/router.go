package app

import (
	"net/http"

	authhandler "github.com/ChernykhITMO/Wishlist-API/internal/auth/handler"
	bookinghandler "github.com/ChernykhITMO/Wishlist-API/internal/booking/handler"
	gifthandler "github.com/ChernykhITMO/Wishlist-API/internal/gifts/handler"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/jwt"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/middleware"
	wishlisthandler "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/handler"
)

func NewRouter(
	tokens *jwt.Manager,
	authHandler *authhandler.Handler,
	wishlistsHandler *wishlisthandler.Handler,
	giftsHandler *gifthandler.Handler,
	bookingHandler *bookinghandler.Handler,
) http.Handler {
	authMiddleware := middleware.Authenticate(tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("POST /wishlists", authMiddleware(http.HandlerFunc(wishlistsHandler.Create)))
	mux.Handle("GET /wishlists", authMiddleware(http.HandlerFunc(wishlistsHandler.List)))
	mux.Handle("GET /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Get)))
	mux.HandleFunc("GET /wishlists/public/{token}", wishlistsHandler.GetPublic)
	mux.Handle("PUT /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Update)))
	mux.Handle("DELETE /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Delete)))
	mux.Handle("POST /wishlists/{wishlistId}/gifts", authMiddleware(http.HandlerFunc(giftsHandler.Create)))
	mux.Handle("GET /wishlists/{wishlistId}/gifts", authMiddleware(http.HandlerFunc(giftsHandler.List)))
	mux.Handle("GET /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Get)))
	mux.Handle("PUT /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Update)))
	mux.Handle("DELETE /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Delete)))
	mux.Handle("POST /wishlists/public/{token}/bookings", http.HandlerFunc(bookingHandler.CreateBooking))

	return mux
}
