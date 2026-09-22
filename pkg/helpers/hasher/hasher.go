package hasher

import "golang.org/x/crypto/bcrypt"

//go:generate mockery --name HashHelper --filename hash_helper.go
type HashHelper interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type Hasher struct{}

func NewHasher() HashHelper {
	return &Hasher{}
}

func (h *Hasher) HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (h *Hasher) CheckPasswordHash(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
