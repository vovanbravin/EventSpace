package main

import (
	"EventSpace/internal/auth"
	"EventSpace/internal/auth/handler"
	"EventSpace/internal/auth/infrastructure"
	"EventSpace/internal/auth/service"
	"EventSpace/internal/auth/storage"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/crypto/bcrypt"
)

func LoadConfig() (*auth.Config, error) {
	doc, err := os.ReadFile("config/auth/config.toml")

	if err != nil {
		return nil, err
	}

	var config auth.Config
	err = toml.Unmarshal(doc, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	config, err := LoadConfig()

	if err != nil {
		log.Fatalf("Error to load auth config: %v", err)
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.Database.Addr)

	if err != nil {
		log.Fatalf("Error connect to db: %v", err.Error())
	}

	defer db.Close()

	userRepo := storage.NewUserRepository(db)
	tokenRepo := storage.NewRefreshTokenRepository(db)
	hasher := infrastructure.NewBcryptHasher(bcrypt.DefaultCost)
	jwtService := infrastructure.NewJWTService(config.JWT.Secret, 15*time.Minute, 20*time.Hour)

	authService := service.NewAuthService(hasher, tokenRepo, userRepo, jwtService)

	router := chi.NewRouter()

	authHandler := handler.NewAuthHandler(authService)

	router.Post("/register", authHandler.RegisterHandler)

	router.Post("/login", authHandler.Login)

	router.Post("/refresh", authHandler.Refresh)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Start auth")
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Error to start api-gateway: %v", err.Error())
	}
}
