package database

import (
	"fmt"
	// "practise/handlers"
	jwt "practise/config"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	jwxt "github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // binding 用於 Gin 的請求體驗證
	Password string `json:"password" binding:"required"`
}

func GenerateToken(username string) (string, error) {
	// 1. 獲取簽名密鑰
	signingKey, err := jwt.GetJWTSigningKey() // 從密鑰管理模組獲取密鑰
	if err != nil {
		return "", fmt.Errorf("failed to get signing key: %v", err)
	}

	// 2. 構建 JWT Token
	// lestrrat-go/jwx 建議使用 Build() 方法，並直接設定 Claims
	tok, err := jwxt.NewBuilder().
		Claim("username", username).                // 設定自定義 Claim
		IssuedAt(time.Now()).                       // 設定簽發時間
		Expiration(time.Now().Add(24 * time.Hour)). // 設定過期時間
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build token: %v", err)
	}

	// 3. 簽名 JWT
	// jwa.HS256 適用於 HMAC 密鑰。如果你用 RSA，這裡要改為 jwa.RS256() 並傳入 RSA 私鑰
	signed, err := jwxt.Sign(tok, jwxt.WithKey(jwa.HS256, signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return string(signed), nil
}

//POST localhost:8080/login
// the string in return is a token
func Login(username, password string) (string,error) {
	//先匹配db
	// collection:= GetCollection("users")
	user,err:=FindUsersByUsername(username)
	if err != nil {
		return "",err
	}
	if user == nil {
		return "",fmt.Errorf("%s is not exists",username)
		//也可以改成跳轉到CreatedUser之類的行為
	}
	// password驗證
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("wrong password")
	}
	token,err := GenerateToken(username)
	if err != nil {
		return "", fmt.Errorf("generate Token fail: %v", err)
	}

	return token,nil
}