package base62_lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Sample struct {
	source      string
	target      string
	sourceBytes []byte
	targetBytes []byte
	alphabet    string
}

func NewSample(source, target string) *Sample {
	return &Sample{source: source, target: target, sourceBytes: []byte(source), targetBytes: []byte(target)}
}

var SamplesStd = []*Sample{
	NewSample("", ""),
	NewSample("f", "1e"),
	NewSample("fo", "6ox"),
	NewSample("foo", "SAPP"),
	NewSample("foob", "1sIyuo"),
	NewSample("fooba", "7kENWa1"),
	NewSample("foobar", "VytN8Wjy"),

	NewSample("su", "7gj"),
	NewSample("sur", "VkRe"),
	NewSample("sure", "275mAn"),
	NewSample("sure.", "8jHquZ4"),
	NewSample("asure.", "UQPPAab8"),
	NewSample("easure.", "26h8PlupSA"),
	NewSample("leasure.", "9IzLUOIY2fe"),

	NewSample("=", "z"),
	NewSample(">", "10"),
	NewSample("?", "11"),
	NewSample("11", "3H7"),
	NewSample("111", "DWfh"),
	NewSample("1111", "tquAL"),
	NewSample("11111", "3icRuhV"),
	NewSample("111111", "FMElG7cn"),

	NewSample("Hello, World!", "1wJfrzvdbtXUOlUjUf"),
	NewSample("你好，世界！", "1ugmIChyMAcCbDRpROpAtpXdp"),
	NewSample("こんにちは", "1fyB0pNlcVqP3tfXZ1FmB"),
	NewSample("안녕하십니까", "1yl6dfHPaO9hroEXU9qFioFhM"),
}

func TestSample(t *testing.T) {

	for i, sample := range SamplesStd {
		t.Run(fmt.Sprintf("sample %d", i), func(t *testing.T) {
			t.Parallel()
			enc := NewStdEncoding()

			// // encode Test
			encodeByte := enc.Encode(sample.sourceBytes)
			assert.Equal(t, sample.targetBytes, encodeByte)

			encodeString := enc.EncodeToString(sample.sourceBytes)
			assert.Equal(t, sample.target, encodeString)

			// decode
			decodeByte, err := enc.Decode(sample.targetBytes)
			assert.NoError(t, err)
			assert.Equal(t, sample.sourceBytes, decodeByte)

			decodeString, err := enc.DecodeString(sample.targetBytes)
			assert.NoError(t, err)
			assert.Equal(t, sample.source, decodeString)
		})
	}
}
