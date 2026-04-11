package app

import (
	"net/http"

	_ "github.com/ChernykhITMO/Wishlist-API/docs"
	authhandler "github.com/ChernykhITMO/Wishlist-API/internal/auth/handler"
	bookinghandler "github.com/ChernykhITMO/Wishlist-API/internal/booking/handler"
	gifthandler "github.com/ChernykhITMO/Wishlist-API/internal/gifts/handler"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/jwt"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/middleware"
	wishlisthandler "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/handler"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	tokens *jwt.Manager,
	authHandler *authhandler.Handler,
	wishlistsHandler *wishlisthandler.Handler,
	giftsHandler *gifthandler.Handler,
	bookingHandler *bookinghandler.Handler,
) http.Handler {
	authMiddleware := middleware.Authenticate(tokens)

	rootMux := http.NewServeMux()
	publicMux := http.NewServeMux()
	privateMux := http.NewServeMux()

	privateMux.HandleFunc("POST /register", authHandler.Register)
	privateMux.HandleFunc("POST /login", authHandler.Login)
	privateMux.Handle("POST /wishlists", authMiddleware(http.HandlerFunc(wishlistsHandler.Create)))
	privateMux.Handle("GET /wishlists", authMiddleware(http.HandlerFunc(wishlistsHandler.List)))
	privateMux.Handle("GET /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Get)))
	privateMux.Handle("PUT /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Update)))
	privateMux.Handle("DELETE /wishlists/{id}", authMiddleware(http.HandlerFunc(wishlistsHandler.Delete)))
	privateMux.Handle("POST /wishlists/{wishlistId}/gifts", authMiddleware(http.HandlerFunc(giftsHandler.Create)))
	privateMux.Handle("GET /wishlists/{wishlistId}/gifts", authMiddleware(http.HandlerFunc(giftsHandler.List)))
	privateMux.Handle("GET /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Get)))
	privateMux.Handle("PUT /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Update)))
	privateMux.Handle("DELETE /wishlists/{wishlistId}/gifts/{giftId}", authMiddleware(http.HandlerFunc(giftsHandler.Delete)))
	privateMux.Handle("GET /swagger/{path...}", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	publicMux.HandleFunc("GET /wishlists/public/{token}", wishlistsHandler.GetPublic)
	publicMux.Handle("POST /wishlists/public/{token}/bookings", http.HandlerFunc(bookingHandler.CreateBooking))

	rootMux.Handle("/wishlists/public/", publicMux)
	rootMux.Handle("/", privateMux)

	return rootMux
}
