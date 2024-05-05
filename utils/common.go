package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}

// HasdPassword returns the bcrypt hash of the password
func HasdPassword(password string) (string, error) {
	hasdedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hasdedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hasdedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasdedPassword), []byte(password))
}
