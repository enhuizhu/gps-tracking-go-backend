package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/enhuizhu/gps-tracking-go-backend/src/models"
	"github.com/enhuizhu/gps-tracking-go-backend/src/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserController to deal with all the related user request
type UserController struct {
}

// Index default method
func (h *UserController) Index() {
	fmt.Println("hello home controller")
}

func getUserDataBaseOnRequest(c *gin.Context) (*models.UserLoginModel, error) {
	rawData, error := c.GetRawData();
	var user models.UserLoginModel
	
	if error != nil {
		return nil, error
	} else {
		json.Unmarshal(rawData, &user)
	}

	return &user, nil;
}

// CreateNewUser for creating new user login
func (u *UserController) CreateNewUser(c *gin.Context) {
	user, error := getUserDataBaseOnRequest(c)

	if error != nil {
		c.JSON(200, gin.H{
			"error": error,
		})
	} else {
		result := user.CreateLogin();

		if result != constants.OK {
			c.JSON(200, gin.H {
				"success": false,
				"message": result,
			})
		} else {
			c.JSON(200, gin.H {
				"success": true,
			})
		}
	}
}

func (u *UserController) Login(c *gin.Context) {
	user, error := getUserDataBaseOnRequest(c)
	// fmt.Println(user)
	if error != nil {
		c.JSON(200, gin.H{
			"error": error,
			// error: "error on parsing data from request",
		})
	} else {
		td, err := user.Login();

		if err != nil {
			c.JSON(200, gin.H{
				"error": err,
			})
		} else {
			tokens := map[string]string{
				"access_token": td.AccessToken,
				"refresh_token": td.RefreshToken,
			}

			c.JSON(200, tokens);
		}
	}
}

func (u *UserController) Logout(c *gin.Context) {
	_, err := models.Logout(c.Request)

	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
		})
	}
}

func (u *UserController) AddFriend(c *gin.Context) {
	mapFriendRequest := map[string][]int{}

	if err := c.ShouldBindJSON(&mapFriendRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return 
	}

	// get user detail
	accessDetails, err := models.ExtractTokenMetadata(c.Request)

	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})

		return 
	}

	userId, err := models.GetUserId(accessDetails.Email)

	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})

		return 
	}

	models.AddFriendRequest(userId, mapFriendRequest["friendIds"])

	c.JSON(200, gin.H{
		"success": true,
	})
}
