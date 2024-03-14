package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MastersAndFans/masterfans-backend/pkg/mocks"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	// Mock environment variable for JWT_SECRET_KEY
	os.Setenv("JWT_SECRET_KEY", "testsecret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	userRepoMock := new(mocks.UserRepoMock)
	authHandler := NewAuthHandler(userRepoMock)

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		t.Fatalf("Hashing mock user password failed: %v", err)
	}
	mockUser := &models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
		ID:       1,
	}

	tests := []struct {
		name           string
		prepare        func()
		input          LoginRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Successful Login",
			prepare: func() {
				userRepoMock.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
			},
			input: LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Invalid Credentials",
			prepare: func() {
				userRepoMock.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Not found!")).Once()
			},
			input: LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()

			body, _ := json.Marshal(tc.input)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			authHandler.LoginHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectToken {
				// Ensure the Set-Cookie header is present and correct.
				cookie := rr.Header().Get("Set-Cookie")
				assert.Contains(t, cookie, "auth_token=")
			} else {
				// Ensure no auth_token cookie is set on failed login attempts.
				cookie := rr.Header().Get("Set-Cookie")
				assert.NotContains(t, cookie, "auth_token=")
			}

			userRepoMock.AssertExpectations(t)
		})
	}
}
