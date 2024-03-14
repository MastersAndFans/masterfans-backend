package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MastersAndFans/masterfans-backend/pkg/auth"
	"github.com/MastersAndFans/masterfans-backend/pkg/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	mockedTime, _ := time.Parse("2006-01-02", "1990-01-01")
	birthDateString := mockedTime.Format("2006-01-02")

	userRepoMock := new(mocks.UserRepoMock)
	authHandler := auth.NewAuthHandler(userRepoMock)

	r := chi.NewRouter()
	r.Post("/api/register", authHandler.RegisterHandler)

	tests := []struct {
		name           string
		prepare        func()
		input          auth.RegisterRequest
		expectedStatus int
	}{
		{
			name: "success",
			prepare: func() {
				userRepoMock.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found")).Once()
				userRepoMock.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()
			},
			input: auth.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123",
				RepeatPass:  "password123",
				Name:        "John",
				Surname:     "Doe",
				BirthDate:   birthDateString,
				PhoneNumber: "123456789",
				IsMaster:    false,
			},
			expectedStatus: http.StatusCreated,
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

			userRepoMock.AssertExpectations(t)
		})
	}
}
