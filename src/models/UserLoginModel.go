package models

import (
	"github.com/enhuizhu/gps-tracking-go-backend/src/db"
	"github.com/enhuizhu/gps-tracking-go-backend/src/constants"
	"github.com/enhuizhu/gps-tracking-go-backend/src/helpers"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"log"
)
// UserLoginModel for dealing user login data
type UserLoginModel struct {
	Email string
	Password string
}

var traceDb = db.Db{}

func hashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost);

	if err != nil {
		log.Println(err)
	}

	return string(hash);
}

func doesPasswordMatch(hashedPassword string, password string) bool {
	pwdBytes := []byte(password);
	hashedBytes := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedBytes, pwdBytes)

	if err != nil {
		return false
	}

	return true
}


func (userLogin *UserLoginModel) CreateLogin() string{
	if !helpers.IsValidEmail(userLogin.Email) {
		return constants.INVALID_EMAIL
	}
	
	var number int
	err := traceDb.QueryRow("select count(*) from user_login where email = ?", userLogin.Email).Scan(&number);
	
	if err != nil {
		panic(err.Error())
	}

	if number > 0 {
		return fmt.Sprintf("user with email %s is already exist.", userLogin.Email)
	}

	traceDb.Query("insert into user_login (email, password) values (?, ?)", userLogin.Email, hashPassword(userLogin.Password))
	
	return constants.OK
}
