package jwt_pkg

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JWTGenerator --filename generator.go
type JwtGenerator interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
}

type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

var JwtAlogorithm = jwt.SigningMethodRS256

// secrete key will get from pem file
func NewJWTGenerator(privateKeyPath string) JwtGenerator {
	privateKeyPem, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPem)
	if err != nil {
		panic(err)
	}

	return &jwtGenerator{privateKey: privateKey}
}

func (g *jwtGenerator) GenerateToken(claims jwt.MapClaims) (string, error) {
	var token = jwt.NewWithClaims(JwtAlogorithm, claims)

	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
