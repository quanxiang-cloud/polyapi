package rule

import (
	"testing"
)

func TestCheckCharSet(t *testing.T) {
	oks := []string{
		"ns2",
		"ns",
		"n s",
		"",
		"\n",
		"中文《》{}【】，。；‘“",
	}
	for i, v := range oks {
		if err := CheckCharSet(v); err != nil {
			panic(i)
		}
	}
}
