package action

import (
	"context"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// CreateNamespace create namespace
func CreateNamespace(parent, name, title, owner, ownerName string, ignoreAccessCheck bool) (string, error) {
	fullPath := apipath.Join(parent, name)
	req := &adaptor.CreateNsReq{
		Namespace:         parent,
		Name:              name,
		Title:             title,
		Owner:             owner,
		OwnerName:         ownerName,
		IgnoreAccessCheck: ignoreAccessCheck,
	}
	if req.Owner == "" {
		req.Owner = consts.SystemName
		req.OwnerName = consts.SystemTitle
	}
	oper := adaptor.GetNamespaceOper()
	if _, err := oper.Create(context.Background(), req); err != nil {
		return fullPath, err
	}
	return fullPath, nil
}
