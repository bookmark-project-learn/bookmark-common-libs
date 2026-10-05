package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt_pkg "github.com/bookmark-project-learn/bookmark-common-libs/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Engine interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	SetupRouteHttpTest(token string) *httptest.ResponseRecorder
}

type engine struct {
	app          *gin.Engine
	jwtValidator jwt_pkg.JwtValidator
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) SetupRouteHttpTest(token string) *httptest.ResponseRecorder {
	// Setup the test route
	e.app.Use(NewJwtAuthMiddleware(e.jwtValidator).JwtAuth())
	e.app.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	// make request recorder
	req, _ := http.NewRequest("GET", "/test", nil)
	tokenHeader := ""
	if token != "" {
		tokenHeader = "Bearer " + token
	}
	req.Header.Set("Authorization", tokenHeader)
	w := httptest.NewRecorder()
	e.app.ServeHTTP(w, req)
	return w
}

func NewEngine(jwtValidator jwt_pkg.JwtValidator) Engine {
	return &engine{
		app:          gin.New(),
		jwtValidator: jwtValidator,
	}
}

const publicKeyPath = "./publickey.test.pem"

func TestJwtAuthMiddleware(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		token      string
		statusCode int
	}{
		{
			name:       "StatusOK",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTc5MTI5NDYxM30.K69mXWM7SQcOq82gJSw5n4f3VrK6bOMRv2SHEVatE1j6oyu7trhEVBne7qIqwyM9Ux4ZCROg7CdenihxXJy179nCClTHR_Us0TsCiHrRbsCGK5yJgFXXP26xSPcI_HeyxUZSe8VRIoMex4BQO1lTJAU29ZkLfGjBfFWpOKPayljIcAaWyv8yM7gpAAZgWeJWgJgpr1WlYps-F2HPD8gXQRtIjKPRLdiwRilf8OJ-JvCVJ_rQmkwMIAoQKeIfYGu45nTOaSPAF-qjTEqP8Dgjk3TahgEvhO5UbZrWjV_s7VIyCAXcBjteZoIQtAWWpCXB7_ov_04Ya-tj4fXWDHT7tw",
			statusCode: http.StatusOK,
		},
		{
			name:       "Not Valid Token",
			token:      "invalid_token",
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "ParseTokenError",
			token:      "Bearer invalid_token_format",
			statusCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			jwtValidator := jwt_pkg.NewJWTValidator(publicKeyPath)
			e := NewEngine(jwtValidator)
			response := e.SetupRouteHttpTest(tc.token)
			assert.Equal(t, tc.statusCode, response.Code, "Expected status code does not match actual status code")
		})
	}
}
