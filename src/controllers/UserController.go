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

// CreateNewUser for creating new user login
func (u * UserController) CreateNewUser(c *gin.Context) {
	rawData, error := c.GetRawData();

	if error != nil {
		c.JSON(200, gin.H{
			"error": error,
		})
	} else {
		var user models.UserLoginModel
		json.Unmarshal(rawData, &user)
		
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