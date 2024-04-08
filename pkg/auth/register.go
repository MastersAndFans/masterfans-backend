package auth

import (
	"encoding/json"
	"github.com/MastersAndFans/masterfans-backend/pkg/helpers"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"net/http"
	"regexp"
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

	if !(isValidEmail(req.Email)) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if len(req.Password) < 8 || !isPasswordStrong(req.Password) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Password is too weak")
		return
	}

	if req.Password != req.RepeatPass {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if !isValidName(req.Name) || !isValidName(req.Surname) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Invalid name or surname")
		return
	}

	if !isValidBirthDate(req.BirthDate) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Invalid birth date")
		return
	}

	_, err := h.config.UserRepo.FindByEmail(r.Context(), req.Email)
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

	if !isValidPhoneNumber(req.PhoneNumber) {
		helpers.ErrorHelper(w, http.StatusBadRequest, "Invalid phone number")
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

	err = h.config.UserRepo.CreateUser(r.Context(), &user)
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

func isValidPhoneNumber(number string) bool {
	var phoneNumberRegex = regexp.MustCompile(`^\+370[0-9]{8,8}$`)
	return phoneNumberRegex.MatchString(number)
}

func isValidBirthDate(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	thirteenYearsAgo := time.Now().AddDate(-13, 0, 0)
	return date.After(thirteenYearsAgo)
}

func isValidName(name string) bool {
	var nameRegex = regexp.MustCompile(`^[a-zA-ZąčęėįšųūžĄČĘĖĮŠŲŪŽ]{2,50}$`)
	return nameRegex.MatchString(name)
}

func isPasswordStrong(password string) bool {
	var (
		minLenRegex      = regexp.MustCompile(`.{8,}`)
		upperCaseRegex   = regexp.MustCompile(`[A-Z]`)
		lowerCaseRegex   = regexp.MustCompile(`[a-z]`)
		numberRegex      = regexp.MustCompile(`[0-9]`)
		specialCharRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)

	return minLenRegex.MatchString(password) &&
		upperCaseRegex.MatchString(password) &&
		lowerCaseRegex.MatchString(password) &&
		numberRegex.MatchString(password) &&
		specialCharRegex.MatchString(password)
}

func isValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
