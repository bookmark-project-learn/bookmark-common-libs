package jwt_pkg

import (
	"crypto/rand"
	"crypto/rsa"
)

type MockJwt struct {
	JwtGenarate JwtGenerator
	JwtValidate JwtValidator
}

func MockPrivateKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.Public().(*rsa.PublicKey)
	return privateKey, publicKey
}

func NewMockJwt() *MockJwt {
	privateKey, publicKey := MockPrivateKey()
	jwtGenarate := &jwtGenerator{privateKey: privateKey}
	jwtValidate := &jwtValidator{publicKey: publicKey}
	return &MockJwt{jwtGenarate, jwtValidate}
}
