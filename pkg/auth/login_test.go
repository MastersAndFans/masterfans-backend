package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MastersAndFans/masterfans-backend/pkg/mocks"
	"github.com/MastersAndFans/masterfans-backend/pkg/models"
	"github.com/MastersAndFans/masterfans-backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "testsecret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	userRepoMock := new(mocks.UserRepoMock)
	authHandler := NewAuthHandler(userRepoMock)

	hashedPassword, err := utils.HashPassword("password")
	assert.NoError(t, err, "Hashing mock user password should not fail")

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
		expectedBody   map[string]string
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
			expectedBody:   map[string]string{"message": "Logged in successfully"},
		},
		{
			name: "Invalid Credentials",
			prepare: func() {
				userRepoMock.On("FindByEmail", mock.Anything, "wrong@example.com").Return(nil, errors.New("not found")).Once()
			},
			input: LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]string{"message": "Invalid credentials"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()

			body, err := json.Marshal(tc.input)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			authHandler.LoginHandler(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Expected JSON response, got parsing error: %v. Response body: %s", err, rr.Body.String())
			}

			assert.Equal(t, tc.expectedBody["message"], responseBody["message"], "Response message does not match expected value")

			userRepoMock.AssertExpectations(t)
		})
	}

}
