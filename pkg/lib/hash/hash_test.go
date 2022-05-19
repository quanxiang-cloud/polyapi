package hash

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	id2 "github.com/quanxiang-cloud/cabin/id"
)

const testShowIDDetail = false

func TestMd5(t *testing.T) {
	h1 := Md5Hash(0, "x", "y")
	h2 := Md5Hash(0, "xy")

	if h1 != h2 {
		fmt.Println("Md5Hash", h1, len(h1))
		fmt.Println("Md5Hash", h2, len(h2))

		fmt.Println(base64Coder.DecodeString(h1))
		fmt.Println(base64Coder.DecodeString(h2))
		t.Fatalf("h1 %s != h2 %s", h1, h2)
	}
}

func TestUUID(t *testing.T) {
	id := id2.HexUUID(false)
	hexid := strings.Replace(id, "-", "", -1)
	b, _ := hex.DecodeString(hexid)
	b64 := base64.RawURLEncoding.EncodeToString(b)
	id64 := GenID("req")
	const hashIdx = 1
	const salt = "salt"
	hex := GenHexID(false)
	hexU := GenHexID(true)

	type testCase struct {
		name   string
		val    string
		expect string
	}
	cases := []*testCase{
		&testCase{
			name:   "UUID",
			val:    id,
			expect: id,
		},
		&testCase{
			name:   "hexUUID",
			val:    hexid,
			expect: hexid,
		},
		&testCase{
			name:   "base64UUID",
			val:    b64,
			expect: b64,
		},
		&testCase{
			name:   "GenID",
			val:    id64,
			expect: id64,
		},
		&testCase{
			name:   "hexid",
			val:    hex,
			expect: hex,
		},
		&testCase{
			name:   "hexidU",
			val:    hexU,
			expect: hexU,
		},
		&testCase{
			name:   "md5",
			val:    Md5Hash(hashIdx, salt),
			expect: "AZ7yk2Nv2pHxSzVxmmpCu8g",
		},
		&testCase{
			name:   "sha1",
			val:    Sha1Hash(hashIdx, salt),
			expect: "AUMHqPpoz-Ly9l1SW7dC3shE10wC",
		},
		&testCase{
			name:   "sha224",
			val:    Sha224Hash(hashIdx, salt),
			expect: "ARXZgF1K3ECzgWY4PaLPpLyn8SJF5mh_1xijRAM",
		},
		&testCase{
			name:   "sha256",
			val:    Sha256Hash(hashIdx, salt),
			expect: "AY_-K44zEci_faIQ7yDeJeaJceci0GfBpWKYTJ2KLow5",
		},
		&testCase{
			name:   "default",
			val:    Default("req", hashIdx, salt),
			expect: "req_AY_-K44zEci_faIQ7yDeJeaJceci0GfBpWKYTJ2KLow5",
		},
	}
	for i, v := range cases {
		if testShowIDDetail {
			fmt.Printf("%-2d %-10s size=%02d  %s\n", i+1, v.name, len(v.val), v.val)
			if "GenID" == v.name {
				fmt.Println()
			}
		}
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d TestUUID(%s) fail: expect %v, got %v", i+1, msg, expect, got)
			}
		}
		assert(v.val, v.expect, v.name)
	}
}
