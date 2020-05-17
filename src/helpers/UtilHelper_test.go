package helpers

import (
	"fmt"
	"testing"
)

func TestEmailValidation(t *testing.T) {
	if !IsValidEmail("afdal@fa.com") {
		t.Error("afdal@fa.com should be valid email")
	}

	if IsValidEmail("afdal") {
		t.Error("afdal is not a valid email")
	}
}

func TestArrayContain(t *testing.T) {
	// var intArr []interface{}
	intArr := []int{1, 2, 3, 4, 5}
	testInt := 2

	if !ArrayContain(intArr, testInt) {
		t.Error("2 should b inside the array")
	}
}

func TestJSONStringify(t *testing.T) {
	intArr := []int{1, 2, 3}
	result, err := JSONStringify(intArr)

	fmt.Print(err)
	fmt.Print(result)
}
