package base62_helper

import (
	"testing"

	base62_lib "github.com/bookmark-common-libs/pkg/lib/base62"
	"github.com/stretchr/testify/assert"
)

func TestBase62(t *testing.T) {
	testCase := []struct {
		name string
		data string
	}{
		{
			name: "normal case",
			data: "1234",
		},
	}

	for _, tc := range testCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			enc := base62_lib.NewStdEncoding()
			base62Helper := New(enc)
			res := base62Helper.Encode(tc.data)
			assert.Equal(t, enc.EncodeToString([]byte(tc.data)), res)

		})

	}
}
