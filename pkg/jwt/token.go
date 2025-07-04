package jwt

import (
	"fmt"
	"practise/global"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

//取得secret

func GenerateToken(username string) (string, error) {
	now := time.Now()
	// claims := new(CustomClaims)
	// claims.Issuer = global.JWTSetting.Issuer
	// claims.IssuedAt = now
	// claims.ExpiresAt = now.Add(global.JWTSetting.Expire)
	// claims.Username = username
	claims := map[string]interface{}{
		"issuer":     global.JWTSetting.Issuer,
		"issued_at":  now,
		"expires_at": now.Add(global.JWTSetting.Expire),
		"username":   username,
	}

	token, err := jwt.NewBuilder().
		Issuer(global.JWTSetting.Issuer).
		IssuedAt(now).
		Expiration(now.Add(global.JWTSetting.Expire)).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build JWT token: %v", err)
	}

	for key, val := range claims {
		err := token.Set(key, val)
		if err != nil {
			return "", fmt.Errorf("<%s> failed to set claim: %s", key, err)
		}
	}

	signedToken, err := jwt.Sign(token,
		jwt.WithKey(global.JWTSetting.Algorithm, []byte(global.JWTSetting.Secret)))
	if err != nil {
		return "", fmt.Errorf("%v failed to sign JWT token", err)
	}
	return string(signedToken), nil
}

func ValidateToken(tokenStr string) (jwt.Token, error) {
	if len(tokenStr) == 0 {
		return nil, fmt.Errorf("token cannot be empty")
	}
	// fmt.Println("tokenStr", tokenStr)  已確定有東西
	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithKey(global.JWTSetting.Algorithm, []byte(global.JWTSetting.Secret)),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate JWT token: %v", err)
	}

	return token, nil
}
