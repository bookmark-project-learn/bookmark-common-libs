package code_gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomCodeGenerator_GenerateRandomCode(t *testing.T) {
	testKeyGen := NewKeyGenerator()
	res := testKeyGen.GenerateRandomCode(7)
	assert.Len(t, res, 7)
	for _, c := range res {
		assert.Contains(t, charset, string(c))
	}
}

func TestGenerateRandomCodeUnique(t *testing.T) {
	res := GenerateRandomCodeUnique(7)
	assert.Len(t, res, 7)
	for _, c := range res {
		assert.Contains(t, charset, string(c))
	}
}
