package models

import (
	// "fmt"
	"testing"
	"log"
)

func testHashPassword(t *testing.T)  {
	hashedPassword := hashPassword("test")
	t.Log("hashedPassword2:" + hashedPassword)
	log.Println("hashedPassword:" + hashedPassword);
}

func testDoesPasswordMatch(t *testing.T)  {
	if !doesPasswordMatch("$2a$04$lIVhDwqYympUCxYQUD9Jde4SwFS8K1kvVTdkrFnTPBEIjvDd4aaO.", "test") {
		t.Error("$2a$04$lIVhDwqYympUCxYQUD9Jde4SwFS8K1kvVTdkrFnTPBEIjvDd4aaO. should match test")
	}
}