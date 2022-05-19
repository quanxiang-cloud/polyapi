package gate

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	fmt.Println(fastSplit("127.0.0.1", '.', 4))
	fmt.Println(fastSplit("127.0.*.1", '.', 4))
	fmt.Println(fastSplit("*.google.com", '.', 4))
}

func TestTireTree(t *testing.T) {
	tt := NewTireTree()
	ips := []string{
		"127.0.0",
		"127.0",
		"127.0.1",
		"127.*.1",
		"127.1.*",
	}
	fmt.Println("init", tt.Show())
	for _, v := range ips {
		fmt.Println("insert", v, tt.Insert(v, black))
		fmt.Println(tt.Match(v))
		fmt.Println(tt.Show())
	}
	fmt.Println("match 127.1.0 :")
	fmt.Println(tt.Match("127.1.0"))
	fmt.Println("match 127.1.1 :")
	fmt.Println(tt.Match("127.1.1"))
	fmt.Println("match 127.0.0 :")
	fmt.Println(tt.Match("127.0.0"))
	for _, v := range ips {
		fmt.Println("----------")
		fmt.Println(tt.Show())
		fmt.Println("delete", v, tt.Delete(v))
		fmt.Println(tt.Match(v))
		fmt.Println(tt.Show())
	}
}
