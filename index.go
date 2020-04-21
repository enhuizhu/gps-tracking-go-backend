package main

import (
	"encoding/json"
	"fmt"

	"github.com/enhuizhu/gps-tracking-go-backend/controllers"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		h := controllers.HomeController{}
		h.Index()
		// dbConnection.CreateCon()
		// CreateCon()
		// sayHello()
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/update", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		rawData, error := c.GetRawData()
		fmt.Println(string(rawData))
		var user User
		json.Unmarshal(rawData, &user)
		fmt.Println(user)
		fmt.Println(user.Username)
		fmt.Println(user.Password)

		if error != nil {
			c.JSON(200, gin.H{
				"error": "there is error when get raw data",
			})
		} else {
			c.JSON(200, gin.H{
				"postData": rawData,
			})
		}
	})

	r.Run()
}
