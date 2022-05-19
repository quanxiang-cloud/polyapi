package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// PathOperator PathOperator
type PathOperator interface {
	UpdateAppValid(ctx context.Context, req *UpdateAppValidReq) (*UpdateAppValidResp, error)
	UpdateValidByPath(ctx context.Context, req *UpdateValidByPathReq) (*UpdateValidByPathResp, error)
	DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error)
	DelPath(ctx context.Context, req *DelPathReq) (*DelPathResp, error)
	ExportApp(ctx context.Context, req *ExportAppReq) (*ExportAppResp, error)
	ExportPath(ctx context.Context, req *ExportPathReq) (*ExportPathResp, error)
	Import(ctx context.Context, req *ImportReq) (*ImportResp, error)
}

// CreateOperator create operator
func CreateOperator(conf *config.Config) (PathOperator, error) {
	return &operator{}, nil
}

type operator struct{}

// UpdateAppValidReq UpdateAppValidReq
type UpdateAppValidReq struct {
	AppID string `json:"-"`
	Valid uint   `json:"valid"`
}

// UpdateAppValidResp UpdateAppValidResp
type UpdateAppValidResp struct {
	Errors []error `json:"errors"`
}

// UpdateAppValid UpdateAppValid
func (opt *operator) UpdateAppValid(ctx context.Context, req *UpdateAppValidReq) (*UpdateAppValidResp, error) {
	_, appID := apipath.Split(req.AppID)
	appPath := app.RootPath(appID)

	errors, err := opt.UpdateValidByPath(ctx, &UpdateValidByPathReq{
		Path:  appPath,
		Valid: req.Valid,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateAppValidResp{
		Errors: errors.Errors,
	}, nil
}

// UpdateValidByPathReq UpdateValidByPathReq
type UpdateValidByPathReq struct {
	Path  string `json:"-"`
	Valid uint   `json:"valid"`
}

// UpdateValidByPathResp UpdateValidByPathResp
type UpdateValidByPathResp struct {
	Errors []error `json:"errors"`
}

// UpdateValidByPath UpdateValidByPath
func (opt *operator) UpdateValidByPath(ctx context.Context, req *UpdateValidByPathReq) (*UpdateValidByPathResp, error) {
	req.Path = apipath.Format(req.Path)
	errors := make([]error, 0)

	if op := adaptor.GetNamespaceOper(); op != nil {
		_, err := op.ValidWithSub(ctx, &adaptor.UpdateNsValidReq{
			NamespacePath: req.Path,
			Valid:         req.Valid,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}

	if op := adaptor.GetRawAPIOper(); op != nil {
		_, err := op.ValidByPrefixPath(ctx, &adaptor.RawValidByPrefixPathReq{
			NamespacePath: req.Path,
			Valid:         req.Valid,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}

	if op := adaptor.GetPolyOper(); op != nil {
		_, err := op.ValidByPrefixPath(ctx, &adaptor.PolyValidByPrefixPathReq{
			NamespacePath: req.Path,
			Valid:         req.Valid,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}

	return &UpdateValidByPathResp{
		Errors: errors,
	}, nil
}

// DelAppReq DelAppReq
type DelAppReq struct {
	AppID string
}

// DelAppResp DelAppResp
type DelAppResp struct {
	Errors []*ErrNode `json:"errors"`
}

// DelApp del app
func (opt *operator) DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error) {
	_, appID := apipath.Split(req.AppID)
	appRoot := app.RootPath(appID)
	errStack := opt.delPath(ctx, appRoot)
	return &DelAppResp{
		Errors: errStack,
	}, nil
}

// DelPathReq DelPathReq
type DelPathReq struct {
	NamespacePath string `json:"-"`
}

// DelPathResp DelPathResp
type DelPathResp struct {
	Errors []*ErrNode `json:"errors"`
}

// ErrNode errNode
type ErrNode struct {
	DB    string `json:"db"`
	Table string `json:"table"`
	SQL   string `json:"sql"`
	Err   error  `json:"error"`
}

// DelPath DelPath
func (opt *operator) DelPath(ctx context.Context, req *DelPathReq) (*DelPathResp, error) {
	errStack := opt.delPath(ctx, req.NamespacePath)
	return &DelPathResp{
		Errors: errStack,
	}, nil
}

func (opt *operator) delPath(ctx context.Context, path string) []*ErrNode {
	errStack := make([]*ErrNode, 0)

	if op := adaptor.GetNamespaceOper(); op != nil {
		_, err := op.InnerDelByPrefixPath(ctx, &adaptor.InnerDelNsByPrefixPathReq{
			NamespacePath: path,
		})
		if err != nil {
			namespace, name := apipath.Split(path)
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "api_namespace",
				SQL:   "delete from api_namespace where parent like '" + path + "%' or (parent='" + namespace + "' and namespace='" + name + "')",
				Err:   err,
			})
		}
	}

	if op := adaptor.GetRawAPIOper(); op != nil {
		_, err := op.InnerDelByPrefixPath(ctx, &adaptor.InnerDelRawByPrefixPathReq{
			NamespacePath: path,
		})
		if err != nil {
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "api_raw",
				SQL:   "delete from api_raw where namespace like '" + path + "%'",
				Err:   err,
			})
		}
	}

	if op := adaptor.GetPolyOper(); op != nil {
		_, err := op.InnerDelByPrefixPath(ctx, &adaptor.InnerDelPolyByPrefixPathReq{
			NamespacePath: path,
		})
		if err != nil {
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "api_poly",
				SQL:   "delete from api_poly where namespace like '" + path + "%'",
				Err:   err,
			})
		}
	}

	if op := adaptor.GetServiceOper(); op != nil {
		_, err := op.InnerDelByPrefixPath(ctx, &adaptor.InnerDelServiceByPrefixPathReq{
			NamespacePath: path,
		})
		if err != nil {
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "api_service",
				SQL:   "delete from api_service where namespace like '" + path + "%'",
				Err:   err,
			})
		}
	}

	if op := adaptor.GetRawPolyOper(); op != nil {
		_, err := op.InnerDelByPrefixPath(ctx, &adaptor.InnerDelRawPolyByPrefixPathReq{
			NamespacePath: path,
		})
		if err != nil {
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "api_raw_poly",
				SQL:   "delete from api_raw_poly where raw_api like '" + path + "%' or poly_api like '" + path + "%'",
				Err:   err,
			})
		}
	}

	if op := adaptor.GetKMSOper(); op != nil {
		_, err := op.DeleteCustomerKeyByPrefix(ctx, &adaptor.DeleteCustomerKeyByPrefixReq{
			NamespacePath: path,
		})
		if err != nil {
			errStack = append(errStack, &ErrNode{
				DB:    "mysql",
				Table: "customer_secret_key",
				SQL:   "delete from customer_secret_key where service like '" + path + "%'",
				Err:   err,
			})
		}
	}

	return errStack
}

// ExportAppReq ExportAppReq
type ExportAppReq struct {
	AppID string `json:"-"`
}

// ExportAppResp ExportAppResp
type ExportAppResp = ExportPathResp

// ExportApp ExportApp
func (opt *operator) ExportApp(ctx context.Context, req *ExportAppReq) (*ExportAppResp, error) {
	_, appID := apipath.Split(req.AppID)
	appPath := app.RootPath(appID)
	return opt.ExportPath(ctx, &ExportPathReq{
		Path: appPath,
	})
}

// ExportPathReq ExportPathReq
type ExportPathReq struct {
	Path string `json:"-"`
}

// ExportPathResp ExportPathResp
type ExportPathResp struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

// ExportData ExportData
type ExportData struct {
	Namespace []*adaptor.APINamespace `json:"ns"`
	Raw       []*adaptor.RawAPIFull   `json:"raw"`
	Poly      []*adaptor.PolyAPIFull  `json:"poly"`
	Service   []*adaptor.APIService   `json:"service"`
	RawPoly   []*adaptor.RawPoly      `json:"rawPoly"`
}

// ExportPath ExportPath
func (opt *operator) ExportPath(ctx context.Context, req *ExportPathReq) (*ExportPathResp, error) {
	data := &ExportData{}
	if op := adaptor.GetNamespaceOper(); op != nil {
		ns, err := op.ListByPrefixPath(ctx, &adaptor.ListNsByPrefixPathReq{
			NamespacePath: req.Path,
			Active:        -1,
			Page:          1,
			PageSize:      -1,
		})
		if err != nil {
			return nil, err
		}
		data.Namespace = ns.List
	}

	if op := adaptor.GetRawAPIOper(); op != nil {
		raw, err := op.ListByPrefixPath(ctx, &adaptor.RawListReq{
			NamespacePath: req.Path,
			Active:        -1,
			Page:          1,
			PageSize:      -1,
		})
		if err != nil {
			return nil, err
		}
		data.Raw = raw.List
	}

	if op := adaptor.GetPolyOper(); op != nil {
		poly, err := op.ListByPrefixPath(ctx, &adaptor.ListPolyByPrefixPathReq{
			NamespacePath: req.Path,
			Active:        -1,
			Page:          1,
			PageSize:      -1,
		})
		if err != nil {
			return nil, err
		}
		data.Poly = poly.List
	}

	if op := adaptor.GetServiceOper(); op != nil {
		service, err := op.ListByPrefixPath(ctx, &adaptor.ListServiceByPrefixReq{
			NamespacePath: req.Path,
			Page:          1,
			PageSize:      -1,
		})
		if err != nil {
			return nil, err
		}
		data.Service = service.List
	}

	if op := adaptor.GetRawPolyOper(); op != nil {
		rawPoly, err := op.ListByPrefixPath(ctx, &adaptor.ListRawPolyByPrefixPathReq{
			Path: req.Path,
		})
		if err != nil {
			return nil, err
		}
		data.RawPoly = rawPoly.List
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	s := unsafeByteString(b)

	return &ExportPathResp{
		Path: req.Path,
		Data: s,
	}, nil
}

// ImportReq ImportReq
type ImportReq struct {
	NewID    string            `json:"newID"`
	OldID    string            `json:"oldID"`
	TableMap map[string]string `json:"tableMap"`
	Data     string            `json:"data"`
}

// ImportResp ImportResp
type ImportResp struct{}

// Import data to copy path
func (opt *operator) Import(ctx context.Context, req *ImportReq) (*ImportResp, error) {
	if req.TableMap == nil {
		req.TableMap = make(map[string]string)
	}

	if _, ok := req.TableMap[req.OldID]; ok {
		return nil, errcode.ErrImportDuplicateID.NewError()
	}
	req.TableMap[req.OldID] = req.NewID
	for k, v := range req.TableMap { // replace all old id to new id
		req.Data = strings.ReplaceAll(req.Data, k, v)
	}

	data := &ExportData{}
	b := unsafeStringBytes(req.Data)
	err := json.Unmarshal(b, data)
	if err != nil {
		return nil, err
	}

	if op := adaptor.GetNamespaceOper(); op != nil {
		_, err = op.InnerImport(ctx, &adaptor.InnerImportNsReq{
			List: data.Namespace,
		})
		if err != nil {
			return nil, err
		}
	}

	if op := adaptor.GetRawAPIOper(); op != nil {
		_, err = op.InnerImport(ctx, &adaptor.InnerImportRawReq{
			List: data.Raw,
		})
		if err != nil {
			return nil, err
		}
	}

	if op := adaptor.GetPolyOper(); op != nil {
		_, err = op.InnerImport(ctx, &adaptor.InnerImportPolyReq{
			List: data.Poly,
		})
		if err != nil {
			return nil, err
		}
	}

	if op := adaptor.GetServiceOper(); op != nil {
		_, err = op.InnerImport(ctx, &adaptor.InnerImportServiceReq{
			List: data.Service,
		})
		if err != nil {
			return nil, err
		}
	}

	if op := adaptor.GetRawPolyOper(); op != nil {
		_, err = op.InnerImport(ctx, &adaptor.InnerImportRawPolyReq{
			List: data.RawPoly,
		})
		if err != nil {
			return nil, err
		}
	}

	return &ImportResp{}, nil
}
