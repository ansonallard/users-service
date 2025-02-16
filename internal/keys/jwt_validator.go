package keys

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

// JWKS represents the JSON Web Key Set structure
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	X5C []string `json:"x5c"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
}

// JWTValidator handles JWT validation using JWKS
type JWTValidator struct {
	jwks JWKS
}

// NewJWTValidator creates a new JWT validator with the provided JWKS
func NewJWTValidator(jwks JWKS) *JWTValidator {
	return &JWTValidator{
		jwks: jwks,
	}
}

// getPublicKeyFromX5C extracts the public key from an x5c entry
func (v *JWTValidator) getPublicKeyFromX5C(x5c string) (*rsa.PublicKey, error) {
	// Create PEM block
	pemKey := "-----BEGIN RSA PUBLIC KEY-----\n"
	// Insert newlines every 64 characters
	remaining := x5c
	for len(remaining) > 64 {
		pemKey += remaining[:64] + "\n"
		remaining = remaining[64:]
	}
	if len(remaining) > 0 {
		pemKey += remaining + "\n"
	}
	pemKey += "-----END RSA PUBLIC KEY-----"

	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pemKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	return pub, nil
}

// ValidateToken validates a JWT token using the JWKS
func (v *JWTValidator) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}

		var matchingKey *JWK
		for _, key := range v.jwks.Keys {
			if key.Kid == kid {
				matchingKey = &key
				break
			}
		}

		if matchingKey == nil {
			return nil, errors.New("no matching key found in JWKS")
		}

		publicKey, err := v.getPublicKeyFromX5C(matchingKey.X5C[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}

// CustomClaims represents your JWT claims structure
type CustomClaims struct {
	jwt.Claims
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ValidateTokenWithClaims validates a token and returns custom claims
func (v *JWTValidator) ValidateTokenWithClaims(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}

		var matchingKey *JWK
		for _, key := range v.jwks.Keys {
			if key.Kid == kid {
				matchingKey = &key
				break
			}
		}

		if matchingKey == nil {
			return nil, errors.New("no matching key found in JWKS")
		}

		publicKey, err := v.getPublicKeyFromX5C(matchingKey.X5C[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token with claims: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
