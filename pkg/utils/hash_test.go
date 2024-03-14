package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "test!pass@word123."
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Hashing the password failed: %v", err)
	}

	if hashedPassword == password {
		t.Errorf("Hashsd password is the same as the original password.")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "test!pass@word123."
	wrongPassword := "wrongPASSword111"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Hashing the password failed: %v", err)
	}

	if !CheckPassword(hashedPassword, password) {
		t.Errorf("Password checking failed for the correct password")
	}

	if CheckPassword(hashedPassword, wrongPassword) {
		t.Errorf("Password checking succeeded for the incorrect password.")
	}
}
