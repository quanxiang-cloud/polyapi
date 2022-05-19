package hash

import (
	"fmt"
	"testing"

	id2 "github.com/quanxiang-cloud/cabin/id"
)

const (
	testNameLen = 8
	maxTest     = 10000_0
)

func testRandID(name string, t *testing.T, fn func(int) string) {
	m := map[string]struct{}{}

	errCnt := 0
	for i := 0; i < maxTest; i++ {
		id := fn(testNameLen)
		if _, ok := m[id]; ok {
			fmt.Printf("%s: duplicate short-id %s at %d/%d\n", name, id, i+1, maxTest)
			errCnt++
		}
		m[id] = struct{}{}
		if i == 0 {
			fmt.Println(id)
		}
	}
	if errCnt > 0 {
		fmt.Printf("************%s: len=%d %d/%d duplicate************\n", name, testNameLen, errCnt, maxTest)
	}
}

func TestRandID(t *testing.T) {
	testRandID("ShortID", t, ShortID)
	testRandID("id2.String", t, id2.String)
}

func TestHShortID(t *testing.T) {
	fmt.Println("HShortID", HShortID(32, 0, "seed1"))
	fmt.Println("HShortID", HShortID(0, 0, "seed"))
	fmt.Println("HShortID", HShortID(10, 0, "seed"))
	fmt.Println("HShortID", HShortID(8, 0, "seed1"))
}
