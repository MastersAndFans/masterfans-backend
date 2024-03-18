package models

import (
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"net/http"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"strconv"
)

type UserHandler struct {
	UserRepo repository.IUserRepository
}

func NewUserHandler(userRepo repository.IUserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (handler *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	id_string := chi.URLParam(r, "id")

	// convert id to int64
	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		http.Error(w, "ID must contain only digits", http.StatusBadRequest)
		return
	}

	user, err := handler.UserRepo.FindById(context.Background(), id)

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
}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	users, err := handler.UserRepo.List(context.Background())

	users_json, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK);
	w.Write(users_json)
}
