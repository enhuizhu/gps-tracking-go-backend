package models

import (
	"github.com/enhuizhu/gps-tracking-go-backend/src/db"
	"github.com/enhuizhu/gps-tracking-go-backend/src/constants"
	"github.com/enhuizhu/gps-tracking-go-backend/src/helpers"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tkanos/gonfig"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"
	"strings"
	"errors"
)

var  client *redis.Client

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	
	if len(dsn) == 0 {
	   dsn = "my-redis:6379"
	}

	client = redis.NewClient(&redis.Options{
	   Addr: dsn, //redis port
	})
	
	_, err := client.Ping().Result()
	
	if err != nil {
	   panic(err)
	}
}


type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// UserLoginModel for dealing user login data
type UserLoginModel struct {
	Email string
	Password string
}

type TokenConfig struct {
	AccessSecret string
	RefreshSecret string
}

type AccessDetails struct {
	AccessUuid string
	RefreshUuid string
	Email string
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

func getTokenConfig() (*TokenConfig, error){
	tokenConfig := TokenConfig{}
	
	dir, err := os.Getwd()
	
	if err != nil {
        return nil, err
	}
	
	err = gonfig.GetConf(dir + "/../../tokenConfig.json", &tokenConfig)
	
	if err != nil {
        return nil, err
	}

	return &tokenConfig, nil
}

func CreateToken(email string) (*TokenDetails, error) {
	//Creating Access Token
	tokenConfig, err := getTokenConfig();

	if err != nil {
		return nil, err
	}
	
	os.Setenv("ACCESS_SECRET", tokenConfig.AccessSecret) //this should be in an env file
	os.Setenv("REFRESH_SECRET", tokenConfig.RefreshSecret)

	td := &TokenDetails{}

	td.AtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()
	
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid;
	atClaims["email"] = email
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	
	if err != nil {
	   return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["email"] = email
	rtClaims["exp"]	= td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return nil, err
	 }

	return td, nil
}

func CreateAuth(email string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)  // converting Unix To UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, email, at.Sub(now)).Err()

	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(td.RefreshUuid, email, rt.Sub(now)).Err()

	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(r *http.Request)(*jwt.Token, error) {
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TokenValid(r *http.Request) error  {
	token, err := VerifyToken(r) 

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func ExtractTokenMetadata(r *http.Request)(*AccessDetails, error) {
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		
		if !ok {
			return nil, errors.New("error on getting accessUuid")
		}

		email, ok := claims["email"].(string)

		if !ok {
			return nil, errors.New("error on getting email")
		}

		refreshUuid, ok := claims["refresh_uuid"].(string)

		if !ok {
			return nil, errors.New("error on getting refresh uuid")
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			RefreshUuid: refreshUuid,
			Email: email,
		}, nil
	}

	return nil, errors.New("unknown error")
}

func FetchAuth(authD *AccessDetails) (string, error) {
  email, err := client.Get(authD.AccessUuid).Result()
  
  if err != nil {
     return "", err
  }
 
  return email, nil
}

func IsAuthorized(c *gin.Context) bool {
	var td TokenDetails;
	
	if err := c.ShouldBindJSON(&td); err != nil {
		return false
	}

	tokenAuth, err := ExtractTokenMetadata(c.Request)

	if err != nil {
		return false
	 }

	 _, err = FetchAuth(tokenAuth)
	 
	 if err != nil {
		return false
	 }
	
	 return true
}

func DeleteAuth(givenUuid string) (int64,error) {
	deleted, err := client.Del(givenUuid).Result()
	
	if err != nil {
	   return 0, err
	}
	
	return deleted, nil
}

func RefreshToken(c *gin.Context) {
	tokenConfig, err := getTokenConfig()

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
	}

	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return 
	}

	refreshToken := mapToken["refresh_token"]

	os.Setenv("REFRESH_SECRET", tokenConfig.RefreshSecret)

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		   return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	 })

	 if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	 }

	  //is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims

	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		accessUuid, ok := claims["access_uuid"].(string) //convert the interface to string
		
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		email, ok := claims["email"].(string)

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		 //Delete the previous Refresh Token
		deleted, delErr := DeleteAuth(refreshUuid)
		
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		deleted, delErr = DeleteAuth(accessUuid)
		
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		
		 //Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(email)
		if  createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		saveErr := CreateAuth(email, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
   

func Logout(r *http.Request) (bool, error){
	accessDetails, err := ExtractTokenMetadata(r)

	if err != nil {
		return false, err
	}

	_, err = DeleteAuth(accessDetails.AccessUuid)

	if err != nil {
		return false, err
	}

	_, err = DeleteAuth(accessDetails.RefreshUuid)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (userLogin *UserLoginModel) Login() (*TokenDetails, error) {
	match, err := userLogin.doesMailAndPasswordMatch()

	if err != nil {
		return nil, err
	}

	if match {
		td, err := CreateToken(userLogin.Email)

		if err != nil {
			return nil, err
		}

		err = CreateAuth(userLogin.Email, td)

		if (err != nil) {
			return nil, err
		}

		return td, nil
	}

	return nil, errors.New("email or pssword is wrong")
}


func (userLogin *UserLoginModel) doesMailAndPasswordMatch() (bool, error) {
	var email string
	var password string
	err := traceDb.QueryRow("select email, password from user_login where email = ?", userLogin.Email).Scan(&email, &password);

	if err != nil {
		return false, err
	}

	// check if password match
	if !doesPasswordMatch(password, userLogin.Password) {
		return false, errors.New("password is wrong")
	}

	return true, nil
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
