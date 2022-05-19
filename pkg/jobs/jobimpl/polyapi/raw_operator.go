package polyapi

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"

	"gorm.io/gorm"
)

var errNotSupport = errors.New("not support")
var errNotFound = errors.New("not found")

func newRawAPIOper(db *gorm.DB) *rawAPIOper {
	p := &rawAPIOper{
		db:   db,
		oper: mysql.NewRawAPIRepo(),
	}
	adaptor.SetRawAPIOper(p)
	return p
}

type rawAPIOper struct {
	db   *gorm.DB
	oper models.RawAPIRepo
}

func (r rawAPIOper) Query(c context.Context, req *adaptor.QueryRawAPIReq) (*adaptor.QueryRawAPIResp, error) {
	var rawAPI *models.RawAPICore
	var err error

	ns, name := apipath.Split(req.APIPath)
	rawAPI, err = r.oper.Get(r.db, ns, name)
	if err != nil {
		return nil, err
	}
	if rawAPI.ID == "" {
		return nil, errNotFound
	}

	return &adaptor.QueryRawAPIResp{
		ID:        rawAPI.ID,
		Content:   rawAPI.Content,
		URL:       rawAPI.URL,
		Method:    rawAPI.Method,
		Namespace: rawAPI.Namespace,
		Owner:     rawAPI.Owner,
		OwnerName: rawAPI.OwnerName,
		Service:   rawAPI.Service,
		Schema:    rawAPI.Schema,
		Host:      rawAPI.Host,
		Active:    rawAPI.Active,
		Valid:     rawAPI.Valid,
		Name:      rawAPI.Name,
		Title:     rawAPI.Title,
		UpdateAt:  rawAPI.UpdateAt,
	}, nil

}
func (r rawAPIOper) List(c context.Context, req *adaptor.RawListReq) (*adaptor.RawListResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) ListInService(c context.Context, req *adaptor.ListInServiceReq) (*adaptor.ListInServiceResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) InnerUpdateRawInBatch(ctx context.Context, req *adaptor.InnerUpdateRawInBatchReq) (*adaptor.InnerUpdateRawInBatchResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) InnerDel(c context.Context, req *adaptor.DelReq) (*adaptor.DelResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) ValidInBatches(ctx context.Context, req *adaptor.RawValidInBatchesReq) (*adaptor.RawValidInBatchesResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) ValidByPrefixPath(ctx context.Context, req *adaptor.RawValidByPrefixPathReq) (*adaptor.RawValidByPrefixPathResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) QueryInBatches(c context.Context, req *adaptor.QueryRawAPIInBatchesReq) (*adaptor.QueryRawAPIInBatchesResp, error) {
	return nil, errNotSupport
}

func (r rawAPIOper) ListByPrefixPath(ctx context.Context, req *adaptor.ListRawByPrefixPathReq) (*adaptor.ListRawByPrefixPathResp, error) {
	return nil, errNotSupport
}

func (r rawAPIOper) InnerImport(ctx context.Context, req *adaptor.InnerImportRawReq) (*adaptor.InnerImportRawResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) InnerDelByPrefixPath(ctx context.Context, req *adaptor.InnerDelRawByPrefixPathReq) (*adaptor.InnerDelRawByPrefixPathResp, error) {
	return nil, errNotSupport
}
func (r rawAPIOper) QueryDoc(ctx context.Context, req *apiprovider.QueryDocReq) (*apiprovider.QueryDocResp, error) {
	return &apiprovider.QueryDocResp{
		Doc: json.RawMessage(`{}`),
	}, nil
}
func (r rawAPIOper) QuerySwagger(ctx context.Context, req *adaptor.QueryRawSwaggerReq) (*adaptor.QueryRawSwaggerResp, error) {
	return nil, errNotSupport
}
