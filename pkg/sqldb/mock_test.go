package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMockMemPostgres(t *testing.T) {
	_, err := NewMiniPostgres(t)
	assert.NoError(t, err)
}
