package example

import "fmt"

func hello() {
	str := "just try something, will not use"
	fmt.Println(str)
}

// import (
// 	"bytes"
// 	"fmt"
// 	"net/http"
// 	"practise/global"
// 	"practise/pkg/setting"
// 	"time"

// 	"github.com/lestrrat-go/jwx/v3/jwa"
// 	"github.com/lestrrat-go/jwx/v3/jwe"
// 	"github.com/lestrrat-go/jwx/v3/jwk"
// 	"github.com/lestrrat-go/jwx/v3/jws"
// 	"github.com/lestrrat-go/jwx/v3/jwt"
// )

// func setupSetting() error {
// 	setting, err := setting.NewSetting()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = setting.ReadSection("Server", &global.JWTSetting)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	return nil
// }

// func Example() {
// 	setupSetting()
// 	jsonRSAPrivateKey := []byte(global.JWTSetting.JsonRSAPrivateKey)
// 	privkey, err := jwk.ParseKey(jsonRSAPrivateKey)
// 	if err != nil {
// 		fmt.Printf("failed to parse JWK: %s\n", err)
// 		return
// 	}

// 	pubkey, err := jwk.PublicKeyOf(privkey)
// 	if err != nil {
// 		fmt.Printf("failed to get public key: %s\n", err)
// 		return
// 	}

// 	// Work with JWTs!
// 	{
// 		// Build a JWT!
// 		tok, err := jwt.NewBuilder().
// 			Issuer(`you_need_a_issuer`).
// 			IssuedAt(time.Now()).
// 			Build()
// 		if err != nil {
// 			fmt.Printf("failed to build token: %s\n", err)
// 			return
// 		}
// 		//claim iss & iat
// 		//typ default JWT

// 		// Sign a JWT!
// 		signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256(), privkey))
// 		if err != nil {
// 			fmt.Printf("failed to sign token: %s\n", err)
// 			return
// 		}
// 		//alg is RS256

// 		/*
// 			now
// 			header
// 			{
// 				alg: RS256
// 				typ: JWT
// 			}
// 			payload
// 			{
// 				iss: `you_need_a_issuer`
// 				iat: time.Now()
// 			}
// 		*/
// 		//Signature會在這時一併計算出來

// 		// Verify a JWT!
// 		{
// 			verifiedToken, err := jwt.Parse(signed, jwt.WithKey(jwa.RS256(), pubkey))
// 			if err != nil {
// 				fmt.Printf("failed to verify JWS: %s\n", err)
// 				return
// 			}
// 			_ = verifiedToken
// 		}

// 		// Work with *http.Request!
// 		{
// 			req, err := http.NewRequest(http.MethodGet, `看是什麼函式調用就寫他的url`, nil)
// 			req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))

// 			verifiedToken, err := jwt.ParseRequest(req, jwt.WithKey(jwa.RS256(), pubkey))
// 			if err != nil {
// 				fmt.Printf("failed to verify token from HTTP request: %s\n", err)
// 				return
// 			}
// 			_ = verifiedToken
// 		}
// 	}

// 	// Encrypt and Decrypt arbitrary payload with JWE!
// 	// JWE用於加密整個JWT或任意數據
// 	// 文檔演示使用非對稱加密來加密及解密數據
// 	{
// 		//func jwe.Encrypt(payload []byte, options ...jwe.EncryptOption) ([]byte, error)
// 		//Encrypt第一個引數payload為「真正要加密的原始數據」，通常會將敏感的個資用這加密
// 		//第二個引數 options 如何加密的選項
// 		encrypted, err := jwe.Encrypt(payloadLoremIpsum, jwe.WithKey(jwa.RSA_OAEP(), jwkRSAPublicKey))
// 		//jwa.RSA_OAEP()： 這是一個 jwa.KeyEncryptionAlgorithm 類型的值，指定了金鑰加密演算法。RSA_OAEP 是 RSA 演算法的一種模式，常用於加密金鑰。
// 		//jwkRSAPublicKey： 這是公鑰 (Public Key)。在非對稱加密（如 RSA）中，使用接收方的公鑰來加密數據，這樣只有擁有對應私鑰的接收方才能解密。
// 		if err != nil {
// 			fmt.Printf("failed to encrypt payload: %s\n", err)
// 			return
// 		}

// 		decrypted, err := jwe.Decrypt(encrypted, jwe.WithKey(jwa.RSA_OAEP(), jwkRSAPrivateKey))
// 		//參考Encrypt
// 		if err != nil {
// 			fmt.Printf("failed to decrypt payload: %s\n", err)
// 			return
// 		}

// 		//report param1 & param2 is same length & same bytes
// 		if !bytes.Equal(decrypted, payloadLoremIpsum) {
// 			fmt.Printf("verified payload did not match\n")
// 			return
// 		}
// 	}

// 	// Sign and Verify arbitrary payload with JWS!
// 	// 任意數據!! 是任意數據!!不一定是JWT
// 	{
// 		//原始數據,加密/解密規則,驗證的公/私鑰
// 		signed, err := jws.Sign(payloadLoremIpsum, jws.WithKey(jwa.RS256(), jwkRSAPrivateKey))
// 		if err != nil {
// 			fmt.Printf("failed to sign payload: %s\n", err)
// 			return
// 		}

// 		verified, err := jws.Verify(signed, jws.WithKey(jwa.RS256(), jwkRSAPublicKey))
// 		if err != nil {
// 			fmt.Printf("failed to verify payload: %s\n", err)
// 			return
// 		}

// 		if !bytes.Equal(verified, payloadLoremIpsum) {
// 			fmt.Printf("verified payload did not match\n")
// 			return
// 		}
// 	}
// 	// OUTPUT:
// }
