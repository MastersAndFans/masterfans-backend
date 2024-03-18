package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogoutHandler(t *testing.T) {
	handler := &AuthHandler{}
	req, err := http.NewRequest("POST", "/api/logout", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(handler.LogoutHandler)

	handlerFunc.ServeHTTP(rr, req)

	// Verify status code
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Verify response body
	expected := `{"message":"Logged out successfully"}`
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")

	// Verify cookie
	cookie := rr.Result().Cookies()
	assert.NotEmpty(t, cookie, "Expected cookie to be set")
	expectedCookieName := "auth_token"
	if assert.NotEmpty(t, cookie, "Expected at least one cookie") {
		assert.Equal(t, expectedCookieName, cookie[0].Name, "Expected cookie name to match")
		assert.Equal(t, "", cookie[0].Value, "Expected cookie value to be empty")
		assert.True(t, cookie[0].Expires.Unix() < time.Now().Unix(), "Expected cookie to be expired")
	}
}
