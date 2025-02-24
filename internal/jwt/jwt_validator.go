package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/ansonallard/users-service/internal/jwk"
	"github.com/golang-jwt/jwt"
)

// JWKS represents the JSON Web Key Set structure

// JWTValidator handles JWT validation using JWKS
type JWTValidator struct {
	jwks jwk.JWKs
}

// NewJWTValidator creates a new JWT validator with the provided JWKS
func NewJWTValidator(jwks jwk.JWKs) *JWTValidator {
	return &JWTValidator{
		jwks: jwks,
	}
}

// getPublicKeyFromX5C extracts the public key from an x5c entry
func (v *JWTValidator) getPublicKeyFromX5C(x5c string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := b64.StdEncoding.DecodeString(x5c)
	if err != nil {
		return nil, err
	}
	publicKeyPem, rest := pem.Decode(publicKeyBytes)
	fmt.Println(rest)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyPem.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// ValidateToken validates a JWT token using the JWKS
func (v *JWTValidator) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return v.validateToken(token)
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}

// ValidateTokenWithClaims validates a token and returns custom claims
func (v *JWTValidator) ValidateTokenWithClaims(tokenString string) (*jwk.CustomClaims, error) {
	claims := &jwk.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return v.validateToken(token)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token with claims: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}

func (v *JWTValidator) validateToken(token *jwt.Token) (*rsa.PublicKey, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid header not found in token")
	}

	var matchingKey *jwk.JWK
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
}
