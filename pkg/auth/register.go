package auth

import (
	"context"
	"encoding/json"
	"github.com/MastersAndFans/masterfans-backend/pkg/helpers"
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
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Email == "" {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Email is required")
		return
	}

	if req.Password == "" {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Password is required")
		return
	}

	if req.Name == "" {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Name is required")
		return
	}

	if req.Surname == "" {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Surname is required")
		return
	}

	if req.BirthDate == "" {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Birth date is required")
		return
	}

	_, err := h.UserRepo.FindByEmail(context.Background(), req.Email)
	if err == nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, "User with this email address already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Invalid birth date format")
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
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := map[string]string{
		"message": "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create JSON")
	}
}
