package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Default parameters for Argon2id
var DefaultParams = &params{
	memory:      64 * 1024, // 64MB
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (hashedPassword string, err error) {
	// Generate random salt
	salt := make([]byte, DefaultParams.saltLength)
	if _, err = rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultParams.iterations,
		DefaultParams.memory,
		DefaultParams.parallelism,
		DefaultParams.keyLength,
	)

	// Encode parameters, salt, and hash into string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash
	hashedPassword = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		DefaultParams.memory,
		DefaultParams.iterations,
		DefaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return hashedPassword, nil
}

// VerifyPassword checks if a password matches a hash
func VerifyPassword(password, encodedHash string) (match bool, err error) {
	// Extract params, salt, and hash from encoded string
	var p *params
	var salt, hash []byte
	p, salt, hash, err = decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the password with the same params and salt
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// Compare hashes in constant time
	match = subtle.ConstantTimeCompare(hash, otherHash) == 1
	return match, nil
}

func decodeHash(encodedHash string) (p *params, salt []byte, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported algorithm")
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}

	p = &params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&p.memory,
		&p.iterations,
		&p.parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
