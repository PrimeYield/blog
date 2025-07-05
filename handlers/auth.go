package handlers

import (
	"fmt"
	"net/http"
	"practise/global"
	"practise/pkg/user"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username string
	Password string
}

func LoginHandler(c *gin.Context) {
	userInfo := new(UserInfo)
	err := c.ShouldBindJSON(userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err.Error()),
		})
		return
	}
	token,err := user.Login(userInfo.Username,userInfo.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	tokenExpiration := time.Now().Add(global.JWTSetting.Expire)
	maxAge := int(time.Until(tokenExpiration).Seconds())
	c.SetCookie(
		"jwt_token",
		token,
		maxAge,
		"/",
		"localhost",
		true,
		true,
	)
	c.JSON(http.StatusOK,gin.H{
		"message": "Login Succeful",
		// "token": token,
	})
}