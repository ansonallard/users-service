package service

import (
	"encoding/base64"
	"errors"

	"github.com/ansonallard/users-service/internal/jwk"
)

type JWKServiceConfig struct {
	PublicKey []byte
}

type JWKService interface {
	GetJWKs() (*jwk.JWKs, error)
}

type jwkService struct {
	publicKey []byte
}

func NewJWKService(config JWKServiceConfig) (JWKService, error) {
	if len(config.PublicKey) == 0 {
		return nil, errors.New("public key not provided")
	}
	return &jwkService{publicKey: config.PublicKey}, nil
}

func (s *jwkService) GetJWKs() (*jwk.JWKs, error) {
	kid := "1"
	x5c := base64.StdEncoding.EncodeToString(s.publicKey)
	jwk := jwk.JWKs{
		Keys: []*jwk.JWK{
			{
				Kid: kid,
				Kty: "RSA",
				Use: "sig",
				Alg: "RS256",
				X5c: []string{x5c},
			},
		},
	}
	return &jwk, nil
}
