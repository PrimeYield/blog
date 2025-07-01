package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

var jwtSecretBytes []byte
var once sync.Once

func LoadJWTSigningKey(secret string) error {
	once.Do(func() {
		if secret == "" {
			log.Println("JWT secret key can not be empty")
			return
		}
		jwtSecretBytes = []byte(secret)
	})
	if jwtSecretBytes == nil {
		return fmt.Errorf("JWT secret key not loaded")
	}
	return nil
}

func GetJWTSigningKey() (jwk.Key, error) {
	if jwtSecretBytes == nil {
		return nil, fmt.Errorf("JWT secret key not initialized. Call LoadJWTSigningKey first")
	}
	// 非對稱演算法要從RSA/ECDSA私鑰中生成JWK
	key, err := jwk.FromRaw(jwtSecretBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK from raw secret: %v", err)
	}
	return key, nil
}
