package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "test!pass@word123."
	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err, "Hashing the password should not produce an error")

	assert.NotEqual(t, password, hashedPassword, "Hashed password should not be the same as the original password")
}

func TestCheckPassword(t *testing.T) {
	password := "test!pass@word123."
	wrongPassword := "wrongPASSword111"

	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err, "Hashing the password should not produce an error")

	assert.True(t, CheckPassword(hashedPassword, password), "Password checking should succeed for the correct password")

	assert.False(t, CheckPassword(hashedPassword, wrongPassword), "Password checking should fail for the incorrect password")
}
