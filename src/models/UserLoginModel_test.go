package models

import (
	// "fmt"
	"testing"
	"log"
	"fmt"
)

func TestHashPassword(t *testing.T)  {
	hashedPassword := hashPassword("test")
	t.Log("hashedPassword2:" + hashedPassword)
	log.Println("hashedPassword:" + hashedPassword);
}

func TestDoesPasswordMatch(t *testing.T)  {
	if !doesPasswordMatch("$2a$04$B0.zUxxFWccKFxX2/xHbau0ls681qrs/oWrtqDCUt/OnwRXrMSMCe", "test") {
		t.Error("$2a$04$B0.zUxxFWccKFxX2/xHbau0ls681qrs/oWrtqDCUt/OnwRXrMSMCe should match test")
	}
}

func TestCreateToken(t *testing.T) {
	userLoginModel := UserLoginModel{
		Email: "test@test.com",
		Password: "test",
	}

	token, err := userLoginModel.createToken()

	if err != nil {
		t.Error("err is:" + err.Error())
	} else {
		fmt.Println(token);
	}
}