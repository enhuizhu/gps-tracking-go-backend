package helpers

import (
	"testing"
)

func TestEmailValidation(t *testing.T)  {
	if !IsValidEmail("afdal@fa.com") {
		t.Error("afdal@fa.com should be valid email")
	}

	if IsValidEmail("afdal") {
		t.Error("afdal is not a valid email")
	}
}