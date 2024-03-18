package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MastersAndFans/masterfans-backend/internal/db"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"encoding/json"
	"strconv"
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

	// get by id
	r.Get("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id_string := chi.URLParam(r, "id")

		// convert id to int64
		id, err := strconv.ParseInt(id_string, 10, 64)
		if err != nil {
			http.Error(w, "ID must contain only digits", http.StatusBadRequest)
			return
		}

		user, err := userRepo.FindById(context.Background(), id)

		if err != nil {
			http.Error(w, fmt.Sprintf("User with ID %d not found", id), http.StatusBadRequest)
			return
		}

		user_json, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK);
		w.Write(user_json)
	})

	// list 
	// TODO
	r.Get("/api/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("list")
		w.WriteHeader(http.StatusOK);
	})

	r.Post("/api/register", authHandler.RegisterHandler)

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
