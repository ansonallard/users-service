package keys

// import (
// 	"crypto/rsa"
// 	"crypto/x509"
// 	"encoding/pem"
// 	"fmt"
// 	"os"
// )

// type JWKJson struct {
// 	Keys *[]JWKKey `json:"keys"`
// }

// type JWKKey struct {
// 	X5c []string `json:"x5c"`
// 	Kty string   `json:"kty"`
// 	Alg string   `json:"alg"`
// 	Use string   `json:"use"`
// 	Kid string   `json:"kid"`
// }

// func encodePublicKey(key *rsa.PublicKey, filename string) ([]byte, error) {
// 	publicKeyBytes := x509.MarshalPKCS1PublicKey(key)
// 	publicKeyBlock := pem.Block{
// 		Type:  "RSA PUBLIC KEY",
// 		Bytes: publicKeyBytes,
// 	}
// 	publicPem, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Printf("error when creating %s %s\n", filename, err)
// 		return nil, err
// 	}
// 	err = pem.Encode(publicPem, &publicKeyBlock)
// 	if err != nil {
// 		fmt.Printf("error when encoding %s %s\n", filename, err)
// 		return nil, err
// 	}
// 	result, _ := os.ReadFile(filename)
// 	return result, nil
// }

// func encodePrivateKey(key *rsa.PrivateKey, filename string) error {
// 	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
// 	privateKeyBlock := pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: privateKeyBytes,
// 	}
// 	privatePem, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Printf("error when creating %s %s\n", filename, err)
// 		return err
// 	}
// 	err = pem.Encode(privatePem, &privateKeyBlock)
// 	if err != nil {
// 		fmt.Printf("error when encoding %s %s\n", filename, err)
// 		return err
// 	}
// 	return nil
// }

// func test() {

// 	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{Id: "1234",
// 	// 	Issuer: "authorization.ansonallard.com", IssuedAt: time.Now().Unix(),
// 	// 	ExpiresAt: time.Now().Add(time.Minute * 5).Unix()})
// 	// signedToken, err := token.SignedString(privateKey)

// 	// if err != nil {
// 	// 	fmt.Printf("%+v", err)
// 	// 	os.Exit(1)
// 	// }
// 	// fmt.Printf("jwt: %s\n", signedToken)

// 	publicKeyFileBytes, err := os.ReadFile("public.pem")
// 	// var publicKey *rsa.PublicKey
// 	if err != nil {
// 		if publicKeyFileBytes, err = encodePublicKey(&privateKey.PublicKey, "public.pem"); err != nil {
// 			os.Exit(1)
// 		}
// 	} else {
// 		publicKeyPem, _ := pem.Decode(publicKeyFileBytes)
// 		_, _ = x509.ParsePKCS1PublicKey(publicKeyPem.Bytes)
// 	}

// 	// mux.HandleFunc("/well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
// 	// 	jwk := JWKJson{
// 	// 		Keys: &[]JWKKey{{Kid: "1", Kty: "RSA", Use: "sig", Alg: "RSA256",
// 	// 			X5c: []string{b64.StdEncoding.EncodeToString(publicKeyFileBytes)}}},
// 	// 	}
// 	// 	w.Header().Set("Content-Type", "application/json")
// 	// 	w.WriteHeader(200)
// 	// 	json.NewEncoder(w).Encode(jwk)
// 	// })

// 	// myToken, err := jwt.ParseWithClaims(signedToken, jwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
// 	// 	// Validate signing method
// 	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 	// 	}
// 	// 	return publicKey, nil
// 	// })

// 	// privateKeyFileBytes, err := os.ReadFile("private.pem")
// 	// var privateKey *rsa.PrivateKey
// 	// if err != nil {
// 	// 	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
// 	// 	if err != nil {
// 	// 		fmt.Printf("Cannot generate RSA Key\n")
// 	// 		os.Exit(1)
// 	// 	}
// 	// 	if err = encodePrivateKey(privateKey, "private.pem"); err != nil {
// 	// 		os.Exit(1)
// 	// 	}
// 	// } else {
// 	// 	privateKeyPem, _ := pem.Decode(privateKeyFileBytes)
// 	// 	privateKey, _ = x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
// 	// }
// }
