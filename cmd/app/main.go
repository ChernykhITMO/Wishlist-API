package main

import (
	"context"
	"github.com/ChernykhITMO/Wishlist-API/internal/app"
	"github.com/ChernykhITMO/Wishlist-API/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os/signal"
	"syscall"
)

// @title Wishlist API
// @version 1.0
// @description REST API для создания и управления вишлистами.
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
