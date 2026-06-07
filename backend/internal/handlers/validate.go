package handlers

import (
	"net/mail"
	"strings"
	"unicode"
)

// validateEmail validates and lowercase-normalizes the email.
func validateEmail(email string) (string, bool) {
	email = strings.TrimSpace(strings.ToLower(email))
	if len(email) > 254 {
		return "", false
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return strings.ToLower(addr.Address), true
}

// validatePassword requires >=8 chars with letters and digits; max length caps Argon2 DoS.
func validatePassword(pw string) bool {
	if len(pw) < 8 || len(pw) > 128 {
		return false
	}
	var hasLetter, hasDigit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
