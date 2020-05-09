package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/enhuizhu/gps-tracking-go-backend/src/models"
	"github.com/enhuizhu/gps-tracking-go-backend/src/constants"
	"github.com/gin-gonic/gin"
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

	if error != nil {
		c.JSON(200, gin.H{
			"error": error,
		})
	} else {
		td, err := user.Login();

		if err != nil {
			c.JSON(200, gin.H{
				"error": error,
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


func (u *UserController) RefreshToken(c *gin.Context) {
	models.RefreshToken(c);
}