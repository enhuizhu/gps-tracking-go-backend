package helpers

import (
	"encoding/json"
	"regexp"

	"github.com/gin-gonic/gin"
)

// IsValidEmail for validating email address
func IsValidEmail(email string) bool {
	validEmail := regexp.MustCompile(`^.+@.+\..+$`)
	return validEmail.MatchString(email)
}

// OutputErr common func
func OutputErr(err error, c *gin.Context) {
	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})

		return
	}
}

//OutputMsg common function to send back response
func OutputMsg(success bool, msg string, c *gin.Context) {
	c.JSON(200, gin.H{
		"success": success,
		"msg":     msg,
	})
}

// ArrayContain check if array container any element
func ArrayContain(arr []int, el interface{}) bool {
	for _, e := range arr {
		if e == el {
			return true
		}
	}

	return false
}

// JSONStringify to return json string
func JSONStringify(input interface{}) (string, error) {
	result, err := json.Marshal(input)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
