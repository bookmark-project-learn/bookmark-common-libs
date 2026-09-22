package base62_lib

import (
	"math"
	"strconv"
)

const bits_length = 6

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // must url-safe

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base62 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// ------ Interface ------
// a binary-to-text encoding based on Base62 algorithm that represents arbitrary data
type Base62Encoding interface {
	Encode(src []byte) []byte
	EncodeToString(src []byte) string
	Decode(src []byte) ([]byte, error)
	DecodeString(src []byte) (string, error)
}

// ------ Struct ------
type base62Encoding struct {
	decodeMap [256]byte // map byte 256 to index 0 - 61
}

// ------ initialize Encoding struct ------
func NewStdEncoding() Base62Encoding {
	enc := &base62Encoding{}
	for i := range enc.decodeMap {
		enc.decodeMap[i] = 255 // default value 0xFF
	}
	for i := range encodeStd {
		enc.decodeMap[encodeStd[i]] = byte(i)
	}
	return enc
}

// ------ implementation ------
// - Encode
func (b *base62Encoding) Encode(src []byte) []byte {
	srcLen := len(src)
	maxLengthOfBase62String := int(math.Ceil(math.Log(256) / math.Log(62) * float64(srcLen)))
	destinarionSlice := make([]byte, maxLengthOfBase62String)
	lengthOfBase62String := 0
	// convert destinarionSlice to byte range on 62 values
	for i := range src {
		valueBase256 := int(src[i]) // 256 because the char is a byte 256
		countRealLenghtOfBase62String := 0
		for j := maxLengthOfBase62String - 1; j >= 0 && (valueBase256 != 0 ||
			countRealLenghtOfBase62String < lengthOfBase62String); j-- {
			valueBase256 += 256 * int(destinarionSlice[j])
			destinarionSlice[j] = byte(valueBase256 % 62)
			valueBase256 /= 62
			countRealLenghtOfBase62String += 1
		}
		lengthOfBase62String = countRealLenghtOfBase62String
	}

	// map destinarionSlice from 62 values scheme to byte 256 scheme (abcd...)
	for i := range destinarionSlice {
		destinarionSlice[i] = encodeStd[destinarionSlice[i]]
	}

	if maxLengthOfBase62String > lengthOfBase62String {
		return destinarionSlice[maxLengthOfBase62String-lengthOfBase62String:]
	}
	return destinarionSlice
}

// -- EncodeToString
func (b *base62Encoding) EncodeToString(src []byte) string {
	return string(b.Encode(src))
}

// -- Decode
func (b *base62Encoding) Decode(src []byte) ([]byte, error) {
	srcLen := len(src)
	if srcLen == 0 {
		return []byte{}, nil
	}

	maxLengthOfBase256String := int(math.Ceil(math.Log(62) / math.Log(256) * float64(srcLen)))
	destinarionSlice := make([]byte, maxLengthOfBase256String)
	lengthOfBase62String := 0
	// convert destinarionSlice to byte range on 62 values
	for i := range src {
		valueBase62 := int(b.decodeMap[src[i]]) // 256 because the char is a byte 256
		if valueBase62 == 255 {                 // corrupt whent meet char not mapped to any index, because it is default settting char 0xFF
			return nil, CorruptInputError(src[i])
		}
		countRealLenghtOfBase62String := 0
		for j := maxLengthOfBase256String - 1; j >= 0 && (valueBase62 != 0 ||
			countRealLenghtOfBase62String < lengthOfBase62String); j-- {
			valueBase62 += 62 * int(destinarionSlice[j])
			destinarionSlice[j] = byte(valueBase62 % 256)
			valueBase62 /= 256
			countRealLenghtOfBase62String += 1
		}
		lengthOfBase62String = countRealLenghtOfBase62String
	}

	if maxLengthOfBase256String > lengthOfBase62String {
		return destinarionSlice[maxLengthOfBase256String-lengthOfBase62String:], nil
	}

	return destinarionSlice, nil
}

// -- DecodeString
func (b *base62Encoding) DecodeString(src []byte) (string, error) {

	decode, err := b.Decode(src)
	if err != nil {
		return "", err
	}
	return string(decode), nil
}

// ------ end implementation ------
