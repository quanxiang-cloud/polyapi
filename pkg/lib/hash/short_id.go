package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const (
	// DefaultShortNameLen is defult length of short name
	DefaultShortNameLen = 12
)
const (
	alphaTab   = "abcdfghjklmnpqrstvwxz012456789"
	trimBits   = 1
	tabLen     = len(alphaTab)
	headTabLen = tabLen - 10 // first byte dont allow number character

)

// HShortID generate random short id with hash
func HShortID(n int, index int, elems ...string) string {
	if n <= 0 {
		n = DefaultShortNameLen
	}

	h := sha256.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
		h.Write([]byte{'\n', 0})
	}
	h.Write(unsafeStringBytes(fmt.Sprintf("index_%d", index)))

	bs := h.BlockSize()
	b := make([]byte, 0, (n+bs-1)/bs*bs)
	for i := 0; ; i++ {
		h.Write(unsafeStringBytes(fmt.Sprintf("batch_%d", i)))
		b = h.Sum(b)
		if len(b) >= n {
			break
		}
	}

	b = b[:n]
	mod := headTabLen
	for i, v := range b {
		idx := (int(v>>trimBits) % mod)

		mod = tabLen

		b[i] = alphaTab[idx]
	}
	return string(b)
}

// ShortID  generate a random string with length n
func ShortID(n int) string {
	s, err := ShortIDWithError(n)
	if err != nil {
		panic(err)
	}
	return s
}

// ShortIDWithError  generate a random string with length n
func ShortIDWithError(n int) (string, error) {
	if n <= 0 {
		n = DefaultShortNameLen
	}
	b := make([]byte, n)
	if nr, err := rand.Read(b); err != nil || nr != len(b) {
		return "", err
	}

	mod := headTabLen
	for i, v := range b {
		idx := (int(v>>trimBits) % mod)
		mod = tabLen

		b[i] = alphaTab[idx]
	}
	return string(b), nil
}
