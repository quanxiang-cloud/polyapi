package value

import (
	cryptorand "crypto/rand"
	"math/rand"
	"unsafe"
)

const (
	defaultMockStringLen = 10
	maxMockNumber        = 20
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-"

// MockString generate a random mock string with random length
func MockString() String {
	size := rand.Int()%defaultMockStringLen + 2
	return String(RandString(defaultCharset, size))
}

// MockNumber generate a random number
func MockNumber() Number {
	return Number(rand.Int() % maxMockNumber)
}

// RandString generate a random string by charset
func RandString(charset string, length int) string {
	if charset == "" {
		charset = defaultCharset
	}
	if length <= 0 {
		length = defaultMockStringLen
	}
	charsetSize := len(charset)
	buf := make([]byte, length)
	cryptorand.Read(buf)
	for i := 0; i < length; i++ {
		buf[i] = charset[int(buf[i])%charsetSize]
	}
	return unsafeByteString(buf)
}

// unsafeByteString convert []byte to string without copy
// the origin []byte **MUST NOT** accessed after that
func unsafeByteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
