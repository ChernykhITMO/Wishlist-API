package app

import (
	"context"
	"fmt"
	authrepo "github.com/ChernykhITMO/Wishlist-API/internal/auth/repository"
	"net/http"
	"time"

	authhandler "github.com/ChernykhITMO/Wishlist-API/internal/auth/handler"
	authservice "github.com/ChernykhITMO/Wishlist-API/internal/auth/services"
	bookinghandler "github.com/ChernykhITMO/Wishlist-API/internal/booking/handler"
	bookingrepo "github.com/ChernykhITMO/Wishlist-API/internal/booking/repository"
	bookingservice "github.com/ChernykhITMO/Wishlist-API/internal/booking/services"
	"github.com/ChernykhITMO/Wishlist-API/internal/config"
	gifthandler "github.com/ChernykhITMO/Wishlist-API/internal/gifts/handler"
	giftrepo "github.com/ChernykhITMO/Wishlist-API/internal/gifts/repository"
	giftservice "github.com/ChernykhITMO/Wishlist-API/internal/gifts/services"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/jwt"
	"github.com/ChernykhITMO/Wishlist-API/internal/platform/password"
	platformpostgres "github.com/ChernykhITMO/Wishlist-API/internal/platform/postgres"
	wishlisthandler "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/handler"
	wishlistrepo "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/repository"
	wishlistservice "github.com/ChernykhITMO/Wishlist-API/internal/wishlists/services"
)

func Run(ctx context.Context, cfg *config.Config) error {
	if err := waitForDatabase(ctx, cfg.DBConfig.URL, "migrations"); err != nil {
		return err
	}

	db, err := platformpostgres.NewPool(ctx, cfg.DBConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	authRepo := authrepo.New(db)
	wishlistsRepo := wishlistrepo.New(db)
	giftsRepo := giftrepo.New(db)
	bookingRepo := bookingrepo.New(db)

	passwords := password.New(12)
	tokens := jwt.New(cfg.JWTSecret, time.Hour*24)

	authSrv := authservice.New(authRepo, tokens, passwords)
	wishlistsSrv := wishlistservice.New(wishlistsRepo)
	giftsSrv := giftservice.New(giftsRepo, giftsRepo)
	bookingSrv := bookingservice.New(bookingRepo)

	authHandler := authhandler.New(*authSrv)
	wishlistsHandler := wishlisthandler.New(*wishlistsSrv, giftsSrv)
	giftsHandler := gifthandler.New(*giftsSrv)
	bookingHandler := bookinghandler.New(bookingSrv)

	server := &http.Server{
		Addr:         cfg.HTTPConfig.Addr,
		Handler:      NewRouter(tokens, authHandler, wishlistsHandler, giftsHandler, bookingHandler),
		ReadTimeout:  cfg.HTTPConfig.ReadTimeout,
		WriteTimeout: cfg.HTTPConfig.WriteTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		fmt.Println("starting application...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case <-ctx.Done():
	case err := <-serverErr:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPConfig.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func waitForDatabase(ctx context.Context, dsn, migrationsDir string) error {
	const retryDelay = 2 * time.Second

	ticker := time.NewTicker(retryDelay)
	defer ticker.Stop()

	for {
		if err := platformpostgres.Migrate(ctx, dsn, migrationsDir); err == nil {
			return nil
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
	}
}
