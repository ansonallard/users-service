package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type JWKJson struct {
	Keys *[]JWKKey `json:"keys"`
}

type JWKKey struct {
	X5c []string `json:"x5c"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
}

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Extract public key from private key
	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func EncodePublicKey(key *rsa.PublicKey, filename string) ([]byte, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(key)
	publicKeyBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error when creating %s %s\n", filename, err)
		return nil, err
	}
	err = pem.Encode(publicPem, &publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encoding %s %s\n", filename, err)
		return nil, err
	}
	result, _ := os.ReadFile(filename)
	return result, nil
}

func EncodePrivateKey(key *rsa.PrivateKey, filename string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error when creating %s %s\n", filename, err)
		return err
	}
	err = pem.Encode(privatePem, &privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encoding %s %s\n", filename, err)
		return err
	}
	return nil
}
