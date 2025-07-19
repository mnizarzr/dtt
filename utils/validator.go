package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword validates password requirements
func IsValidPassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)

	// Check for at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

// IsValidName validates name format
func IsValidName(name string) bool {
	// Trim spaces and check if not empty
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < 2 || len(trimmed) > 100 {
		return false
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	return nameRegex.MatchString(trimmed)
}
