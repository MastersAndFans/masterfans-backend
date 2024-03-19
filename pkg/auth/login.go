package auth

import (
	"context"
	"encoding/json"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/helpers"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	UserRepo repository.IUserRepository
}

type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthHandler(userRepo repository.IUserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserRepo.FindByEmail(context.Background(), req.Email)
	if err != nil {
		helpers.ErrorHelper(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		helpers.ErrorHelper(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	claims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "masterfans",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET_KEY env variable is not set!")
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	response := map[string]string{
		"message": "Logged in successfully",
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to create JSON")
	}

}
