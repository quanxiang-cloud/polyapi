package exprx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"

	"github.com/quanxiang-cloud/cabin/logger"
)

// field refer
var (
	// a | a[ 0][1] | [0]
	// dont match: 123 | 123.4
	exprArrayIndex = `(?:\[\s*\d+\s*\])`
	exprIdent      = fmt.Sprintf(`(?:(?:[_a-zA-Z]\w*%s*)|%s)`, exprArrayIndex, exprArrayIndex)
	allRef         = consts.PolyFieldAccessAllData // $
	allRefDot      = allRef + "."                  // $.
	allDot         = polyAllDataVarName + "."
	// a.b.c
	// $.a.b.c
	// $.a
	// $a
	// $a.[0].b
	// $a[0].b
	exprFieldRefer = fmt.Sprintf(`(?P<HEAD>^|[^\.])(?P<EXPR>(?:(?P<REF>\$+\.?)|(?:%s\.))(?:%s\.)*%s)`,
		exprIdent, exprIdent, exprIdent)

	regexpField = regexp.MustCompile(exprFieldRefer)
	//a.[0] => a[0]
	regexpArrayIndex = regexp.MustCompile(`(?P<HEAD>\.)\[\s*(?P<INDEX>\d+)\s*\]`)
)

// req1.x.y -> d.req1.x.y
// $.req1 -> d.req1
// $req1 -> d.req1
// $$req1 -> d.req1
func fixFullFieldRefer(expr string) (string, error) {
	rep := regexpField.ReplaceAllStringFunc(expr, func(src string) string {
		elems := regexpField.FindAllStringSubmatch(src, 1)[0]
		head, expr, allRef := elems[1], elems[2], elems[3]

		if allRef != "" {
			expr = strings.Replace(expr, allRef, "", 1)
		}
		fixedExpr := fixArrayIndex(expr)

		return fmt.Sprintf("%s%s.%s", head, polyAllDataVarName, fixedExpr)
	})

	// TODO: check if is valid expr

	return rep, nil
}

//a.[0] => a[0]
func fixArrayIndex(expr string) string {
	ret := regexpArrayIndex.ReplaceAllStringFunc(expr, func(src string) string {
		elems := regexpArrayIndex.FindAllStringSubmatch(src, 1)[0]
		_, index := elems[1], elems[2]
		return fmt.Sprintf("[%s]", index)
	})
	return ret
}

// ValDirectExpr represents an direct JS expression
type ValDirectExpr ValString // direct expression string

// TypeName returns name of the type
func (v ValDirectExpr) TypeName() string { return ExprTypeDirectExpr.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValDirectExpr) DelayedJSONDecode() error {
	return v.Validate() // validate the value
}

// String convert v to string
func (v ValDirectExpr) String() string {
	return string(v)
}

// SetString set a string to Value
func (v *ValDirectExpr) SetString(s string) {
	v.SetStringWithError(s)
}

// SetStringWithError set a string to Value with format check
func (v *ValDirectExpr) SetStringWithError(s string) error {
	_, err := v.checkAndConvertExpr(v.String(), nil)
	if err == nil {
		*v = ValDirectExpr(s)
	}
	return err
}

// ToScript returns the script of this element represent
func (v ValDirectExpr) ToScript(depth int, e Evaler) (string, error) {
	expr, err := v.checkAndConvertExpr(v.String(), e)
	return expr, err
}

// Validate verify the value of object
func (v ValDirectExpr) Validate() error {
	_, err := v.checkAndConvertExpr(v.String(), nil)
	return err
}

func (v ValDirectExpr) checkAndConvertExpr(expr string, e Evaler) (string, error) {
	expr, err := fixFullFieldRefer(expr)
	if err != nil {
		return "", err
	}
	//BUG: //pdCreateNS('/system/app',d.start.appID,'应用')
	// Eval fail? "SyntaxError: SyntaxError: (anonymous): Line 1:32 Unexpected string (and 2 more errors)
	const comment = "//##"
	if strings.HasPrefix(expr, comment) { // comment expr, dont check
		expr = strings.TrimPrefix(expr, comment)
	} else {
		if e != nil {
			if err := e.Eval(expr); err != nil {
				logger.Logger.Warnf("[expr-eval-fail] expr=%q err=%q", expr, err.Error())
				return "", errcode.ErrVMEvalFail.FmtError(expr)
			}
		}
	}

	return expr, nil
}
