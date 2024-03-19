package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserRepo repository.IUserRepository
}

func NewUserHandler(userRepo repository.IUserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (handler *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.UserRepo.List(context.Background())

	users_json, err := json.Marshal(users)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(users_json)
}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id_string := chi.URLParam(r, "id")

	// convert id to int64
	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID must contain only digits"})
		return
	}

	user, err := handler.UserRepo.FindById(context.Background(), id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("User with ID %d not found", id)})
		return
	}

	user_json, err := json.Marshal(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(user_json)
}
