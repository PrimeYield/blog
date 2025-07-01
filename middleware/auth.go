package middleware

import (
	"log"
	"net/http"
	jwt "practise/config"
	"strings"

	"github.com/gin-gonic/gin"
	jwxt "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/lestrrat-go/jwx/v3/jwa"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		// "Bearer" 是一種 HTTP 身份驗證方案，用於在發送請求時將 Token 傳遞給伺服器。
		// 這是業界廣泛接受的標準，表示後面的字串是一個「持有者 Token」。
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			return
		}

		tokenString := parts[1]

		// 1. 獲取驗證密鑰
		// 對於非對稱演算法，這裡應該是「公鑰」  要再去補一下什麼時候公鑰解密私鑰加密，什麼時候私鑰解密公鑰加密
		verificationKey, err := jwt.GetJWTSigningKey()
		if err != nil {
			log.Printf("AuthMiddleware: Failed to get verification key: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		// 2. 解析和驗證 Token
		// jwt.Parse 會自動處理時間戳驗證 (ExpiresAt, IssuedAt, NotBefore)
		tok, err := jwxt.Parse([]byte(tokenString), jwxt.WithKey(jwa.HS256(), verificationKey))
		if err != nil {
			log.Printf("AuthMiddleware: Token validation failed: %v", err)
			if err.Error() == "token is expired" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			return
		}

		// 3. 從 Token 中提取 Claims
		// jwx 允許直接透過 Claim() 方法來獲取 Claim 值
		username, ok := tok.Get("username")
		if !ok || username == nil {
			log.Printf("AuthMiddleware: Username claim missing or invalid in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims: username missing"})
			return
		}

		// 將用戶名存入 Gin Context
		c.Set("username", username.(string))
		c.Next()
	}
}
