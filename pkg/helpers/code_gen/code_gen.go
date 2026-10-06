package code_gen

import (
	"bytes"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ012345689"

// KeyGenerator represents the interface to generate random string keys
//
// KeyGenerator is an interface to generate random string keys
//
//go:generate mockery --name KeyGenerator --filename string.go
type KeyGenerator interface {
	// GenerateRandomCode generates a random string
	GenerateRandomCode(length int) string
}

type randomCodeGenerator struct {
	rng *rand.Rand
}

// NewKeyGenerator returns a KeyGenerator
func NewKeyGenerator() KeyGenerator {
	return &randomCodeGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateKey generates a random string with no special characters
func (r *randomCodeGenerator) GenerateRandomCode(length int) string {
	return randomCode(r.rng, length)
}

// GenerateRandomCodeUnique generates a unique random string
func GenerateRandomCodeUnique(length int) string {
	// Seed the random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randomCode(rng, length)
}

func randomCode(rng *rand.Rand, length int) string {
	var strBuilder bytes.Buffer
	for i := 0; i < length; i++ {
		strBuilder.WriteByte(charset[rng.Intn(len(charset))])
	}
	return strBuilder.String()
}
