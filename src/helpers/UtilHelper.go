package helpers

import (
	"regexp"
)

// isValidEmail for validating email address
func IsValidEmail(email string) bool {
	validEmail := regexp.MustCompile(`^.+@.+\..+$`)
	return validEmail.MatchString(email)
}