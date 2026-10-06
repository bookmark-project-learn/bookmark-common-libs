package middleware

import (
	"net/http"
	"strings"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type JwtAuthMiddleware struct {
	jwtTokenService jwt_pkg.JwtValidator
}

func NewJwtAuthMiddleware(jwtTokenService jwt_pkg.JwtValidator) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		jwtTokenService: jwtTokenService,
	}
}

// JwtAuth: Check Bearer and token parser
func (j *JwtAuthMiddleware) JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			return
		}
		token := parts[1]
		tokenClaims, err := j.jwtTokenService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Parse token error."})
			return
		}
		c.Set("claims", tokenClaims)
		c.Next()
	}
}
