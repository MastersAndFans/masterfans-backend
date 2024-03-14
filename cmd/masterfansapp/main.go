package main

import (
	"github.com/MastersAndFans/masterfans-backend/internal/db"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance.AutoMigrate(&models.User{})

	userRepo := repository.NewUserRepository(dbInstance)

	authHandler := auth.NewAuthHandler(userRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MasterFans!"))
	})
	r.Post("/api/login", authHandler.LoginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware)
		r.Get("/api/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, MasterFans!"))
		})
	})

	log.Fatal(http.ListenAndServe(":5000", r))

}
