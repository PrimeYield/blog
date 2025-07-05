package middleware

import (
	"fmt"
	"log"
	"net/http"
	"practise/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return authMiddleware
}

func authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) == 0{
		c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error": "Authorization header required",
		})
		log.Printf("Authorization header required")
		return
	}

	parts := strings.Split(authHeader," ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error": "Authorization header format must be Bearer <token>",
		})
		log.Print("Authorization header format must be Bearer <token>")
		return
	}

	tokenStr := parts[1]  //存前端傳來的Token  通常就是從Header提取  比較要記的是Authorization: Bearer <token>這個固定的鍵值對

	token,err := jwt.ValidateToken(tokenStr)
	log.Printf("%v",token)
	if err != nil || token == nil {
		fmt.Printf("tokenStr is a wrong request %v",err)
		return
	}
	// username, err := token.Get("username")
	username, ok := token.Get("username")
		if !ok || username == nil {
			log.Println("AuthMiddleware: Username not found in token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token (username missing)"})
			return
		}
	c.Set("username",c.Request.URL.User.Username())
	c.Next()
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		// "Bearer" 是一種 HTTP 身份驗證方案，用於在發送請求時將 Token 傳遞給伺服器。
// 		// 這是業界廣泛接受的標準，表示後面的字串是一個「持有者 Token」。
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
// 			return
// 		}

// 		tokenString := parts[1]

// 		// 1. 獲取驗證密鑰
// 		// 對於非對稱演算法，這裡應該是「公鑰」  要再去補一下什麼時候公鑰解密私鑰加密，什麼時候私鑰解密公鑰加密
// 		verificationKey, err := jwt.GetJWTSigningKey()
// 		if err != nil {
// 			log.Printf("AuthMiddleware: Failed to get verification key: %v", err)
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
// 			return
// 		}

// 		// 2. 解析和驗證 Token
// 		// jwt.Parse 會自動處理時間戳驗證 (ExpiresAt, IssuedAt, NotBefore)
// 		tok, err := jwxt.Parse([]byte(tokenString), jwxt.WithKey(, verificationKey))
// 		// if err != nil {
// 		// 	log.Printf("AuthMiddleware: Token validation failed: %v", err)
// 		// 	if err.Error() == "token is expired" {
// 		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
// 		// 	} else {
// 		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		// 		log.Println(err)
// 		// 	}
// 		// 	return
// 		// }
// 		if err != nil {
// 			log.Printf("AuthMiddleware: Token validation failed: %v", err)

// 			// 使用 errors.Is 判斷是否為特定的錯誤類型
// 			if errors.Is(err, jwxt.ErrTokenExpired()) {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
// 				log.Println(err) // 打印詳細錯誤日誌方便調試
// 			} else {
// 				// 對於其他類型的錯誤，統一返回 Invalid token
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 				log.Println(err) // 打印詳細錯誤日誌方便調試
// 			}
// 			return
// 		}

// 		// 3. 從 Token 中提取 Claims
// 		// jwx 允許直接透過 Claim() 方法來獲取 Claim 值
// 		username, ok := tok.Get("username")
// 		if !ok || username == nil {
// 			log.Printf("AuthMiddleware: Username claim missing or invalid in token")
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims: username missing"})
// 			return
// 		}

// 		// 將用戶名存入 Gin Context
// 		c.Set("username", username.(string))
// 		c.Next()
// 	}
// }
