package main

import (
	"github.com/enhuizhu/gps-tracking-go-backend/src/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/user/createUser", func(c *gin.Context) {
		userController := controllers.UserController{}
		userController.CreateNewUser(c);
	})

	r.POST("/user/login", func(c *gin.Context) {
		userController := controllers.UserController{};
		userController.Login(c);
	})

	r.POST("/user/logout", func(c *gin.Context) {
		userController := controllers.UserController{};
		userController.Logout(c);
	})

	r.GET("/user/refreshToken", func(c *gin.Context) {

	}) 

	r.Run()
}
