package main

import (
	"github.com/MastersAndFans/masterfans-backend/pkg/handlers"
	"log"
	"net/http"

	"github.com/MastersAndFans/masterfans-backend/internal/db"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance.AutoMigrate(&models.User{})

	userRepo := repository.NewUserRepository(dbInstance)

	authHandler := auth.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MasterFans!"))
	})

	r.Get("/api/user/{id}", userHandler.GetUserById)

	r.Get("/api/user", userHandler.ListUsers)

	r.Post("/api/register", authHandler.RegisterHandler)

	r.Post("/api/login", authHandler.LoginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware)
		r.Post("/api/logout", authHandler.LogoutHandler)
	})

	log.Fatal(http.ListenAndServe(":5000", r))

}
