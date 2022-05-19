package expr

import (
	"fmt"
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// regexURLPara is a patten that enable access parameter in API path
var regexURLPara = regexp.MustCompile(`(?sm:(?:/:(?P<P1>\w+)/??)|(?:/\{(?P<P2>\w+)\}/??))`)

// getStringArg get a string arg of given name
func (s ValueSet) getStringArg(name string) (string, bool) {
	for i := 0; i < len(s); i++ {
		if p := &s[i]; p.Name == name && p.In.IsPath() {
			if p.Field != "" {
				return FullVarName(p.Field.String()), true
			}
			if s, ok := p.Data.D.(Stringer); ok {
				return fmt.Sprintf(`"%s"`, s.String()), true
			}
			break // break loop
		}
	}
	return "", false // get parameter fail
}

// "foo" => foo
func trimHT(s string) string {
	if len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

// ReplacePathArgs replace args in path from the input values
func (s *ValueSet) ReplacePathArgs(srcURL string, name string) (string, error) {
	_fmt, _args, err := s.ResolvePathArgs(srcURL, name)
	if err != nil {
		return "", err
	}
	var args = make([]interface{}, 0, len(_args))
	for _, v := range _args {
		args = append(args, trimHT(v))
	}
	return fmt.Sprintf(_fmt, args...), nil
}

// ResolvePathArgs repalce the path with fmt string and return it's fmt args.
// eg: "/api/:x" => ("/api/%v", ["$x"])
func (s *ValueSet) ResolvePathArgs(srcURL string, nodeName string) (string, []string, error) {
	var err error
	paras := []string(nil)
	repURL := regexURLPara.ReplaceAllStringFunc(srcURL, func(src string) string {
		elems := regexURLPara.FindAllStringSubmatch(src, -1)[0]
		p := ""
		for _, v := range elems[1:] {
			if v != "" {
				p = v
				break
			}
		}

		para, ok := "", false
		if p != "" {
			para, ok = s.getStringArg(p)
		}
		if !ok {
			err = errcode.ErrMissingPathArg.FmtError(nodeName, srcURL, p)
			return src
		}
		paras = append(paras, para)
		return "/%v"
	})

	return repURL, paras, err
}

func replacePathArgsUniversal(srcURL string, name string, inputs map[string]interface{}) (string, error) {
	if inputs == nil {
		return "", fmt.Errorf("error: node %s, inputs is not map[string]interface{}", name)
	}

	var err error
	paras := []interface{}(nil)
	repURL := regexURLPara.ReplaceAllStringFunc(srcURL, func(src string) string {
		elems := regexURLPara.FindAllStringSubmatch(src, -1)[0]
		p := ""
		for _, v := range elems[1:] {
			if v != "" {
				p = v
				break
			}
		}

		para, ok := interface{}(nil), false
		if p != "" {
			para, ok = inputs[p]
		}
		if !ok {
			err = errcode.ErrMissingPathArg.FmtError(name, srcURL, p)
			return src
		}
		if _, ok := para.(string); ok {
			paras = append(paras, para)
		}

		return "/%v"
	})
	return fmt.Sprintf(repURL, paras...), err
}
