package hash

import (
	"github.com/google/uuid"
)

// GenID generate uuid by base64 encoding
func GenID(prefix string) string {
	id := base64Coder.EncodeToString(genUUID())
	if prefix != "" {
		id = prefix + "_" + id
	}
	return id
}

// // GenUUID generate uuid string
// func GenUUID() string {
// 	return gUUID().String()
// }

// GenHexID generate uuid by hex encoding
func GenHexID(upper bool) string {
	id := encodeToHexString(genUUID(), upper)
	return id
}

// genUUID generate google/uuid
func genUUID() []byte {
	b, err := gUUID().MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

var gUUID = uuid.New

const hextable = "0123456789abcdef"
const hextableU = "0123456789ABCDEF"

func encodeToHexString(src []byte, upper bool) string {
	dst := make([]byte, len(src)*2)
	encodeHex(dst, src, upper)
	return string(dst)
}

func encodeHex(dst, src []byte, upper bool) int {
	tb := hextable
	if upper {
		tb = hextableU
	}

	for i, v := range src {
		j := (i << 1)
		dst[j] = tb[v>>4]
		dst[j+1] = tb[v&0x0f]
	}
	return len(src) * 2
}
