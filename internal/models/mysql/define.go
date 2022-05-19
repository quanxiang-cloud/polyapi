package mysql

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

var errNotFound = errcode.ErrNotFound.NewError()

var errNotImplement = errcode.Errorf("not implement")
