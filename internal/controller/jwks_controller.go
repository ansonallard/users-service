package controller

import (
	"errors"
	"net/http"

	"github.com/ansonallard/users-service/internal/service"
	"github.com/gin-gonic/gin"
)

type JWKController interface {
	GetJWKs(g *gin.Context) error
}

type JWKConfig struct {
	Service service.JWKService
}

type jwkController struct {
	service service.JWKService
}

func NewJWKController(config *JWKConfig) (JWKController, error) {
	if config.Service == nil {
		return nil, errors.New("service not provided")
	}
	return &jwkController{service: config.Service}, nil
}

func (c *jwkController) GetJWKs(g *gin.Context) error {
	jwk, err := c.service.GetJWKs()
	if err != nil {
		return err
	}
	g.Header("Content-Type", "application/json")
	g.JSON(http.StatusOK, jwk)
	return nil
	// token := jwt.NewWithClaims(
	// 	jwt.SigningMethodRS256,
	// 	jwt.StandardClaims{
	// 		Id:        "1234",
	// 		Issuer:    "authorization.ansonallard.com",
	// 		IssuedAt:  time.Now().Unix(),
	// 		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
	// 	},
	// )
	// token.Header["kid"] = kid
	// privateKeyFileBytes, err := os.ReadFile(constants.PRIVATE_KEY)
	// if err != nil {
	// 	panic("Cannot read file")
	// }
	// var privateKey *rsa.PrivateKey
	// if err != nil {
	// 	// privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	// 	// if err != nil {
	// 	// 	fmt.Printf("Cannot generate RSA Key\n")
	// 	// 	os.Exit(1)
	// 	// }
	// 	// if err = encodePrivateKey(privateKey, "private.pem"); err != nil {
	// 	// 	os.Exit(1)
	// 	// }
	// } else {
	// 	privateKeyPem, _ := pem.Decode(privateKeyFileBytes)
	// 	privateKey, _ = x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	// }
	// signedToken, err := token.SignedString(privateKey)

	// if err != nil {
	// 	fmt.Printf("%+v", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("jwt: %s\n", signedToken)

	// var jwks keys.JWKS
	// jwksJSON, _ := json.Marshal(jwk)
	// if err := json.Unmarshal([]byte(jwksJSON), &jwks); err != nil {
	// 	log.Fatalf("Failed to parse JWKS: %v", err)
	// }

	// validator := keys.NewJWTValidator(jwks)

	// // Validate a token
	// err = validator.ValidateToken(signedToken)
	// if err != nil {
	// 	log.Fatalf("Token validation failed: %v", err)
	// }

	// fmt.Printf("Token is valid.")
	// return nil
}
