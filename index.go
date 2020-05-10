package main

import (
	"github.com/enhuizhu/gps-tracking-go-backend/src/controllers"
	"github.com/gin-gonic/gin"
	"github.com/enhuizhu/gps-tracking-go-backend/src/models"
	"github.com/enhuizhu/gps-tracking-go-backend/src/middlewares"
)

var userController = controllers.UserController{}

func main() {
	r := gin.Default()

	r.POST("/user/createUser", func(c *gin.Context) {
		userController.CreateNewUser(c);
	})

	r.POST("/user/login", func(c *gin.Context) {
		userController.Login(c);
	})

	r.GET("/user/logout", middlewares.Authorized() , func(c *gin.Context) {
		userController.Logout(c);
	})

	r.POST("/user/refreshToken", func(c *gin.Context) {
		models.RefreshToken(c)
	})
	
	r.POST("/user/addFriend", middlewares.Authorized(), func(c *gin.Context) {
		userController.AddFriend(c)
	})

	r.Run()
}
