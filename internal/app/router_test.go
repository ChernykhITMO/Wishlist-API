package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authhandler "github.com/ChernykhITMO/Wishlist-API/internal/auth/handler"
	bookinghandler "github.com/ChernykhITMO/Wishlist-API/internal/booking/handler"
	gifthandler "github.com/ChernykhITMO/Wishlist-API/internal/gifts/handler"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/jwt"
	wishlisthandler "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/handler"
)

func TestNewRouterRoutesPublicWishlistWithoutServeMuxConflict(t *testing.T) {
	router := NewRouter(
		jwt.New("secret", 24*time.Hour),
		&authhandler.Handler{},
		&wishlisthandler.Handler{},
		&gifthandler.Handler{},
		&bookinghandler.Handler{},
	)

	req := httptest.NewRequest(http.MethodGet, "/wishlists/public/test-token", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
