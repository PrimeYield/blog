package jwt

type CustomClaims struct {
	// 注意：jwx 推薦使用 'jwt.Token' 接口來訪問標準 Claims，
	// 這裡直接定義，jwx會自動映射。
	Issuer    string `json:"iss,omitempty"`
	Subject   string `json:"sub,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	JWTID     string `json:"jti,omitempty"`

	Username string `json:"username"`
}
