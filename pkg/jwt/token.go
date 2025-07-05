package jwt

import (
	"fmt"
	"practise/global"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

//取得secret

func GenerateToken(username string) (string,error) {
	now := time.Now()
	claims := new(CustomClaims)
	claims.Issuer = global.JWTSetting.Issuer
	claims.IssuedAt = now
	claims.ExpiresAt = now.Add(global.JWTSetting.Expire)
	claims.Username = username

	token ,err := jwt.NewBuilder().
	Issuer(global.JWTSetting.Issuer).
	IssuedAt(now).
	Expiration(now.Add(global.JWTSetting.Expire)).
	// Claim("username",username).
	Subject(username).
	Build()
	if err != nil {
		return "", fmt.Errorf("failed to build JWT token: %v",err)
	}
	err = token.Set("issuer",claims.Issuer)
	if err != nil {
		 return "", fmt.Errorf("%v failed to set claim: %s",err,claims.Issuer)
	}
	err = token.Set("issued_at",claims.IssuedAt)
	if err != nil {
		 return "", fmt.Errorf("%v failed to set claim: %v",err,claims.IssuedAt)
	}
	err = token.Set("expires_at",claims.ExpiresAt)
	if err != nil {
		 return "", fmt.Errorf("%v failed to set claim: %v",err,claims.ExpiresAt)
	}
	err = token.Set("username",claims.Username)
	if err != nil {
		 return "", fmt.Errorf("%v failed to set claim: %s",err,claims.Username)
	}
	signedToken, err:= jwt.Sign(token,jwt.WithKey(global.JWTSetting.Algorithm,[]byte(global.JWTSetting.Secret)))
	if err != nil {
 		return "", fmt.Errorf("%v failed to sign JWT token",err)
	}
	return string(signedToken),nil
}

func ValidateToken(tokenStr string) (jwt.Token,error){
	if len(tokenStr) == 0 {
		return nil ,fmt.Errorf("token cannot be empty")
	}
	token ,err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithKey(global.JWTSetting.Algorithm,global.JWTSetting.Secret),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate JWT token: %v",err)
	}
	
	return token ,nil
}