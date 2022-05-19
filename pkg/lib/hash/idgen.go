package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"unsafe"
)

// MaxHashConflict presents max time for hash conflict retry
const MaxHashConflict = 2

var base64Coder = base64.RawURLEncoding

// Default generate default hash, if index>0, it try to avoid hash conflict by salt
func Default(prefex string, index int, elems ...string) string {
	h := Sha256Hash(index, elems...)
	if prefex != "" {
		h = prefex + "_" + h
	}
	return h
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Sha1Hash generate sha1 hash, if index>0, it try to avoid hash conflict by salt
func Sha1Hash(index int, elems ...string) string {
	h := sha1.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	buf := make([]byte, h.Size()+1)
	buf[0] = byte(index)
	return base64Coder.EncodeToString(h.Sum(buf[:1]))
}

// Sha224Hash generate sha224 hash, if index>0, it try to avoid hash conflict by salt
func Sha224Hash(index int, elems ...string) string {
	h := sha256.New224()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	buf := make([]byte, h.Size()+1)
	buf[0] = byte(index)
	return base64Coder.EncodeToString(h.Sum(buf[:1]))
}

// Sha256Hash generate sha256 hash, if index>0, it try to avoid hash conflict by salt
func Sha256Hash(index int, elems ...string) string {
	h := sha256.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	buf := make([]byte, h.Size()+1)
	buf[0] = byte(index)
	return base64Coder.EncodeToString(h.Sum(buf[:1]))
}

// Md5Hash generate sha256 hash, if index>0, it try to avoid hash conflict by salt
func Md5Hash(index int, elems ...string) string {
	h := md5.New()
	for _, v := range elems {
		h.Write(unsafeStringBytes(v))
	}
	if index > 0 {
		h.Write(unsafeStringBytes(fmt.Sprintf("salt_%d", index)))
	}
	buf := make([]byte, h.Size()+1)
	buf[0] = byte(index)
	return base64Coder.EncodeToString(h.Sum(buf[:1]))
}
