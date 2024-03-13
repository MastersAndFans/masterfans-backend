package main

import (
	"github.com/MastersAndFans/masterfans-backend/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	// Sample connection, do not leave it here, only connect to db where needed.
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbInstance.AutoMigrate()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Respond to the root route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MasterFans!"))
	})

	// Sample API route
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, MasterFans!"))
		})
	})

	http.ListenAndServe(":5000", r)

}
