package jwt_pkg

import (
	"crypto/rsa"
	"errors"
	"os"
	"reflect"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JwtValidator --filename validator.go
type JwtValidator interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

func NewJWTValidator(publicKeyPath string) JwtValidator {
	publicKeyPem, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPem)
	if err != nil {
		panic(err)
	}

	return &jwtValidator{publicKey: publicKey}
}

var (
	InvalidToken      = errors.New("invalid token")
	ErrorExtractToken = errors.New("error extract token")
)

func (v *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	refAlgo := reflect.TypeOf(JwtAlogorithm)
	tkn, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ensure the token algorithm matches the expected signing method
		if reflect.TypeOf(token.Method) != refAlgo {
			return nil, jwt.ErrSignatureInvalid
		}
		return v.publicKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, InvalidToken
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, ErrorExtractToken

}
