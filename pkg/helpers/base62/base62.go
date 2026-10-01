package base62_helper

import (
	base62_lib "github.com/bookmark-project-learn/bookmark-common-libs/pkg/lib/base62"
)

//go:generate mockery --name=Base62Helper --filename=base62.go --outpkg=mocks
type Base62Helper interface {
	Encode(data string) string
}

type base62Helper struct {
	base62Encoding base62_lib.Base62Encoding
}

func (b *base62Helper) Encode(data string) string {
	res := b.base62Encoding.EncodeToString([]byte(data))
	return res
}

func New(base62Encoding base62_lib.Base62Encoding) Base62Helper {
	return &base62Helper{base62Encoding}
}
