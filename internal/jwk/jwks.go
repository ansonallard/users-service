package jwk

import "github.com/golang-jwt/jwt"

type JWKs struct {
	Keys []*JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	X5c []string `json:"x5c"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
}

// CustomClaims represents your JWT claims structure
type CustomClaims struct {
	jwt.Claims
	Username string `json:"username"`
	Role     string `json:"role"`
}
