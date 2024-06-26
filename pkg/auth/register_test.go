package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	userRepoMock := new(mocks.UserRepoMock)
	authHandlerConfig := auth.AuthHandlerConfig{UserRepo: userRepoMock, JWTSecretKey: "", TokenDuration: 0}
	authHandler := auth.NewAuthHandler(authHandlerConfig)

	r := chi.NewRouter()
	r.Post("/api/register", authHandler.RegisterHandler)

	tests := []struct {
		name           string
		prepare        func()
		input          auth.RegisterRequest
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "success",
			prepare: func() {
				userRepoMock.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("Not found!")).Once()
				userRepoMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
			},
			input: auth.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123A@",
				RepeatPass:  "password123A@",
				Name:        "John",
				Surname:     "Doe",
				BirthDate:   "2019-04-05",
				PhoneNumber: "+37069536785",
				IsMaster:    false,
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"message": "User created successfully",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()

			body, _ := json.Marshal(tc.input)
			request, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, request)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var responseBody map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			for key, expectedValue := range tc.expectedBody {
				assert.Equal(t, expectedValue, responseBody[key], "Expected JSON response does not match for key: %s", key)
			}

			userRepoMock.AssertExpectations(t)
		})
	}
}
