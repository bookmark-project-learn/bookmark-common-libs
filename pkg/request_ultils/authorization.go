package request_ultils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrorExtractClaims = errors.New("error extract token")
	ClaimsNotFound     = errors.New("Claims extract token")
)

func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenClaims, exists := c.Get("claims")

	if !exists {
		return nil, ClaimsNotFound
	}

	if _, ok := tokenClaims.(jwt.MapClaims); !ok {
		return nil, ErrorExtractClaims
	}
	return tokenClaims.(jwt.MapClaims), nil
}

func GetSubjectFromClaims(c *gin.Context) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return string(""), err
	}
	userId, err := claims.GetSubject()
	if err != nil {
		return string(""), err
	}
	return userId, nil
}
