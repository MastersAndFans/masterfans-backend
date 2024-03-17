package auth

import (
	"context"
	"encoding/json"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	RepeatPass  string `json:"repeat_pass"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	BirthDate   string `json:"birth_date"`
	PhoneNumber string `json:"phone_number"`
	IsMaster    bool   `json:"is_master"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if req.Surname == "" {
		http.Error(w, "Surname is required", http.StatusBadRequest)
		return
	}

	if req.BirthDate == "" {
		http.Error(w, "Birth date is required", http.StatusBadRequest)
		return
	}

	_, err := h.UserRepo.FindByEmail(context.Background(), req.Email)
	if err == nil {
		http.Error(w, "User with this email address already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		http.Error(w, "Invalid birth date format", http.StatusBadRequest)
		return
	}

	user := models.User{
		Email:       req.Email,
		Password:    hashedPassword,
		Name:        req.Name,
		Surname:     req.Surname,
		BirthDate:   birthDate,
		PhoneNumber: req.PhoneNumber,
		IsMaster:    req.IsMaster,
	}

	err = h.UserRepo.CreateUser(context.Background(), &user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}
