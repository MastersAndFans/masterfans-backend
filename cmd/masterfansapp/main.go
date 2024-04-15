package main

import (
	"github.com/MastersAndFans/masterfans-backend/internal/db"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/handlers"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = dbInstance.AutoMigrate(&models.User{}, &models.Review{}, &models.Schedule{}, &models.Service{}, &models.TimeRange{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := repository.NewUserRepository(dbInstance)
	auctionRepo := repository.NewAuctionRepository(dbInstance)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	authHandlerConfig := auth.AuthHandlerConfig{UserRepo: userRepo, JWTSecretKey: secretKey, TokenDuration: 24 * time.Hour}
	authHandler := auth.NewAuthHandler(authHandlerConfig)
	userHandler := handlers.NewUserHandler(userRepo)
	auctionHandler := handlers.NewAuctionHandler(auctionRepo)

	r := setupRouter(authHandler, userHandler, auctionHandler)

	log.Println("Starting server on :5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func setupRouter(authHandler *auth.AuthHandler, userHandler *handlers.UserHandler, auctionHandler *handlers.AuctionHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MasterFans!"))
	})

	// User routes
	r.Route("/api/user", func(r chi.Router) {
		r.Get("/{id}", userHandler.GetUserById)
		r.Get("/", userHandler.ListUsers)
	})

	r.Route("/api/auction", func(r chi.Router) {
		r.Get("/{id}", auctionHandler.GetAuctionById)
		r.Get("/", auctionHandler.ListAuctions)
	})

	// Authentication routes
	r.Post("/api/register", authHandler.RegisterHandler)
	r.Post("/api/login", authHandler.LoginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware)
		r.Post("/api/logout", authHandler.LogoutHandler)
	})

	return r
}
