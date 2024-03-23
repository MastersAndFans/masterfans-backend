package auth

import (
	"encoding/json"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/MastersAndFans/masterfans-backend/pkg/helpers"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type AuthHandlerConfig struct {
	UserRepo      repository.IUserRepository
	JWTSecretKey  string
	TokenDuration time.Duration
}

type AuthHandler struct {
	config AuthHandlerConfig
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthHandler(config AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{config: config}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.ErrorHelper(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.config.UserRepo.FindByEmail(r.Context(), req.Email)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "masterfans",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.config.JWTSecretKey))
	if err != nil {
		helpers.ErrorHelper(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(h.config.TokenDuration),
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
