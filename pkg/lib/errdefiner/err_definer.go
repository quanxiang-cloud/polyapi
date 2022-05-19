package errdefiner

import (
	"fmt"
	"strings"

	error2 "github.com/quanxiang-cloud/cabin/error"
)

// NOTE: create only one ErrorDefiner within an application on initialize
var inst *ErrorDefiner
var _ = NewErrorDefiner() // initialize inst

// NewErrorDefiner create an ErrorDefiner or return inst
func NewErrorDefiner() *ErrorDefiner {
	if inst == nil {
		p := &ErrorDefiner{
			codeTable: make(map[int64]string),
			cacheErr:  make(map[ErrorCode]*error2.Error),
		}
		error2.CodeTable = p.codeTable
		inst = p
	}
	return inst
}

// ErrorDefiner object
type ErrorDefiner struct {
	codeTable map[int64]string            //  error code list
	cacheErr  map[ErrorCode]*error2.Error //cache fix error
}

// MustReg regist an error with duplicate code check
func (r *ErrorDefiner) MustReg(code int64, msg string) ErrorCode {
	if _, ok := r.codeTable[code]; ok {
		panic(fmt.Errorf("duplicate code %d", code))
	}
	r.codeTable[code] = msg
	c := ErrorCode(code)
	if strings.IndexByte(msg, '%') < 0 { // without format parmeter
		_ = c.NewError() // generate cacheErr
	}
	return c
}

func (r *ErrorDefiner) newError(c ErrorCode) *error2.Error {
	if cache, ok := r.cacheErr[c]; ok {
		return cache
	}

	err := error2.New(c.Int64())
	r.cacheErr[c] = &err // cache fix error
	return &err
}

func (r *ErrorDefiner) msg(c ErrorCode, paras []interface{}) string {
	if m, ok := r.codeTable[c.Int64()]; ok {
		return fmt.Sprintf(m, paras...)
	}
	return fmt.Sprintf("<unknown error code %d>", c)
}

// ErrorCode of int
type ErrorCode int64

// Int64 convert code to int64
func (c ErrorCode) Int64() int64 {
	return int64(c)
}

// NewError create an error without format
func (c ErrorCode) NewError() *error2.Error {
	return inst.newError(c)
}

// Msg format an error message string
func (c ErrorCode) Msg(paras ...interface{}) string {
	return inst.msg(c, paras)
}

// FmtError create an error with format
func (c ErrorCode) FmtError(paras ...interface{}) error {
	err := error2.New(c.Int64(), paras...)
	return &err
}

//------------------------------------------------------------------------------

// exports
const (
	ErrParams = error2.ErrParams
	Internal  = error2.Internal
	Unknown   = error2.Unknown
	Success   = error2.Success
)

// Errorf format an standard error with parameters
func Errorf(format string, paras ...interface{}) error {
	return fmt.Errorf(format, paras...)
}

// NewErrorWithString return an error with message
func NewErrorWithString(code int64, msg string) *error2.Error {
	err := error2.NewErrorWithString(code, msg)
	return &err
}
