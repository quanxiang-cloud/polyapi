package expr

import (
	"testing"
)

func TestUrlPara(t *testing.T) {
	url := `https://api.xxx.com/api/v1/:xx/yy/{zz}`
	expect := `https://api.xxx.com/api/v1/p1/yy/p2`
	repURL := regexURLPara.ReplaceAllStringFunc(url, func(src string) string {
		elem := regexURLPara.FindAllStringSubmatch(src, -1)[0]
		p1, p2 := elem[1], elem[2]
		// println(src, p1, p2)
		switch {
		case p1 != "":
			return "/p1"
		case p2 != "":
			return "/p2"
		}
		return ""
	})
	if repURL != expect {
		t.Errorf("expect %s got %s", expect, repURL)
	}
}
