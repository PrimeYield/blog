package handlers

import (
	"fmt"
	"net/http"
	"practise/database"
	"practise/global"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	loginRequest := new(LoginRequest)
	if err := c.ShouldBindJSON(loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err.Error()),
		})
		return
	}
	token, err := database.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	tokenExpiration := time.Now().Add(global.JWTSetting.Expire)
	maxAge := int(time.Until(tokenExpiration).Seconds())
	c.SetCookie(
		"jwt_token",
		token,
		// int(expiration.Sub(time.Now()).Seconds()),
		maxAge,
		"/", // Cookie 的路徑，"/" 表示整個網站都可用
		"localhost",
		true, // Secure: 只在 HTTPS 連接時發送。生產環境強烈建議設為 true。
		true, // HttpOnly: 設為 true，表示 JavaScript 無法訪問此 Cookie。這是關鍵！
	)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
