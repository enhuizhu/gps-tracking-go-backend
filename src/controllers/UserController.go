package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/enhuizhu/gps-tracking-go-backend/src/constants"
	"github.com/enhuizhu/gps-tracking-go-backend/src/helpers"
	"github.com/enhuizhu/gps-tracking-go-backend/src/models"
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
	rawData, error := c.GetRawData()
	var user models.UserLoginModel

	if error != nil {
		return nil, error
	} else {
		json.Unmarshal(rawData, &user)
	}

	return &user, nil
}

// CreateNewUser for creating new user login
func (u *UserController) CreateNewUser(c *gin.Context) {
	user, error := getUserDataBaseOnRequest(c)

	if error != nil {
		helpers.OutputErr(error, c)
	} else {
		result := user.CreateLogin()

		if result != constants.OK {
			c.JSON(200, gin.H{
				"success": false,
				"message": result,
			})
		} else {
			c.JSON(200, gin.H{
				"success": true,
			})
		}
	}
}

func (u *UserController) Login(c *gin.Context) {
	user, error := getUserDataBaseOnRequest(c)
	// fmt.Println(user)
	if error != nil {
		helpers.OutputErr(error, c)
	} else {
		td, err := user.Login()

		if err != nil {
			helpers.OutputErr(err, c)
		} else {
			tokens := map[string]string{
				"access_token":  td.AccessToken,
				"refresh_token": td.RefreshToken,
			}

			c.JSON(200, tokens)
		}
	}
}

func (u *UserController) Logout(c *gin.Context) {
	_, err := models.Logout(c.Request)

	if err != nil {
		helpers.OutputErr(err, c)
	} else {
		c.JSON(200, gin.H{
			"success": true,
		})
	}
}

func (u *UserController) AddFriend(c *gin.Context) {
	mapFriendRequest := map[string]string{}

	if err := c.ShouldBindJSON(&mapFriendRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// get user detail
	accessDetails, err := models.ExtractTokenMetadata(c.Request)

	if err != nil {
		helpers.OutputErr(err, c)
		return
	}

	userID, err := models.GetUserId(accessDetails.Email)

	if err != nil {
		helpers.OutputErr(err, c)
		return
	}

	friendID, err := models.GetUserId(mapFriendRequest["email"])

	if err != nil {
		helpers.OutputMsg(false, mapFriendRequest["email"]+" is not in our system yet, please ask them to join it first", c)
		return
	}

	// should check if they are already friends
	if models.AreTheyFriends(userID, friendID) {
		helpers.OutputMsg(false, mapFriendRequest["email"]+" is already your friend", c)
		return
	}

	if models.AddFriendRequest(userID, friendID) {
		helpers.OutputMsg(true, "Friend request has been send to you friend successfully!", c)
	} else {
		helpers.OutputMsg(false, "same friend request is already there", c)
	}
}

func (u *UserController) AcceptUserRequest(c *gin.Context) {
	requestID := c.Param("requestID")
	// need to check if it's valid request id
	accessDetail, err := models.ExtractTokenMetadata(c.Request)

	if err != nil {
		helpers.OutputErr(err, c)
		return
	}

	userID, err := models.GetUserId(accessDetail.Email)

	if err != nil {
		helpers.OutputErr(err, c)
		return
	}

	requestIDInt, err := strconv.Atoi(requestID)

	if err != nil {
		helpers.OutputErr(err, c)
		return
	}

	_, friendID := models.GetUserIDAndFriendIDBaseOnRequestID(requestIDInt)

	if friendID != userID {
		helpers.OutputMsg(false, "you are not allowed to accept this request", c)
		return
	}

	if models.AcceptFriendRequest(requestIDInt) {
		helpers.OutputMsg(true, "accept user request successfully", c)
	} else {
		helpers.OutputMsg(false, "some error happend when accept the user request", c)
	}
}
