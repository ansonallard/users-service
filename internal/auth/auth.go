package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

type HashParams struct {
	Memory      uint32 `bson:"memory"`
	Iterations  uint32 `bson:"iterations"`
	Parallelism uint8  `bson:"parallelism"`
	SaltLength  uint32 `bson:"salt_length"`
	KeyLength   uint32 `bson:"key_length"`
}

// Default parameters for Argon2id
var DefaultParams = &HashParams{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (hashedPasswordParams *EncodedHashParams, err error) {
	// Generate random salt
	salt := make([]byte, DefaultParams.SaltLength)
	if _, err = rand.Read(salt); err != nil {
		return nil, err
	}

	// Hash the password
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultParams.Iterations,
		DefaultParams.Memory,
		DefaultParams.Parallelism,
		DefaultParams.KeyLength,
	)

	// Encode parameters, salt, and hash into string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	hashedPassword := EncodedHashParams{
		Argon2Version: argon2.Version,
		B64Salt:       b64Salt,
		B64Hash:       b64Hash,
		HashParams: HashParams{
			Memory:      DefaultParams.Memory,
			Iterations:  DefaultParams.Iterations,
			Parallelism: DefaultParams.Parallelism,
			SaltLength:  DefaultParams.SaltLength,
			KeyLength:   DefaultParams.KeyLength,
		},
	}

	return &hashedPassword, nil
}

type EncodedHashParams struct {
	HashParams
	Argon2Version int    `bson:"argon2_version"`
	B64Salt       string `bson:"b64_salt"`
	B64Hash       string `bson:"b64_hash"`
}

// VerifyPassword checks if a password matches a hash
func VerifyPassword(password string, encodedHash *EncodedHashParams) (match bool, err error) {
	salt, err := base64.RawStdEncoding.DecodeString(encodedHash.B64Salt)
	if err != nil {
		return false, err
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(encodedHash.B64Hash)
	if err != nil {
		return false, err
	}

	// Hash the password with the same params and salt
	incomingHash := argon2.IDKey(
		[]byte(password),
		salt,
		encodedHash.Iterations,
		encodedHash.Memory,
		encodedHash.Parallelism,
		encodedHash.KeyLength,
	)

	// Compare hashes in constant time
	match = subtle.ConstantTimeCompare(storedHash, incomingHash) == 1
	return match, nil
}
