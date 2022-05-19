package service

import (
	"errors"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	error2 "github.com/quanxiang-cloud/cabin/error"
)

var errNotFound = error2.NewErrorWithString(error2.ErrParams, "not found")

var errNotImplement = errors.New("not implement")

var errNotSupport = errors.New("not support")

// Operation exports
type Operation = adaptor.Operation

// operations
const (
	OpCreate     = adaptor.OpCreate
	OpUpdate     = adaptor.OpUpdate
	OpEdit       = adaptor.OpEdit
	OpDelete     = adaptor.OpDelete
	OpQuery      = adaptor.OpQuery
	OpBuild      = adaptor.OpBuild
	OpRequest    = adaptor.OpRequest
	OpAddRawAPI  = adaptor.OpAddRawAPI
	OpAddPolyAPI = adaptor.OpAddPolyAPI
	OpAddService = adaptor.OpAddService
	OpAddSub     = adaptor.OpAddSub
	OpPublish    = adaptor.OpPublish
)
