package service

import (
	"context"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"gorm.io/gorm"
)

// APIPermit is the api permit operator
type APIPermit interface {
	CheckAPIPermit(c context.Context, req *CheckPermitReq) (*CheckPermitResp, error)

	AddPermitGroup(c context.Context, req *AddPermitGroupReq) (*AddPermitGroupResp, error)
	DelPermitGroup(c context.Context, req *DelPermitGroupReq) (*DelPermitGroupResp, error)
	UpdatePermitGroup(c context.Context, req *UpdatePermitGroupReq) (*UpdatePermitGroupResp, error)
	ActivePermitGroup(c context.Context, req *ActivePermitGroupReq) (*ActivePermitGroupResp, error)
	ListPermitGroup(c context.Context, req *ListPermitGroupReq) (*ListPermitGroupResp, error)
	QueryPermitGroup(c context.Context, req *QueryPermitGroupReq) (*QueryPermitGroupResp, error)

	AddPermitElem(c context.Context, req *AddPermitElemReq) (*AddPermitElemResp, error)
	UpdatePermitElem(c context.Context, req *UpdatePermitElemReq) (*UpdatePermitElemResp, error)
	DelPermitElem(c context.Context, req *DelPermitElemReq) (*DelPermitElemResp, error)
	ActivePermitElem(c context.Context, req *ActivePermitElemReq) (*ActivePermitElemResp, error)
	ListPermitElem(c context.Context, req *ListPermitElemReq) (*ListPermitElemResp, error)
	QueryPermitElem(c context.Context, req *QueryPermitElemReq) (*QueryPermitElemResp, error)

	AddPermitGrant(c context.Context, req *AddPermitGrantReq) (*AddPermitGrantResp, error)
	DelPermitGrant(c context.Context, req *DelPermitGrantReq) (*DelPermitGrantResp, error)
	UpdatePermitGrant(c context.Context, req *UpdatePermitGrantReq) (*UpdatePermitGrantResp, error)
	ActivePermitGrant(c context.Context, req *ActivePermitGrantReq) (*ActivePermitGrantResp, error)
	ListPermitGrant(c context.Context, req *ListPermitGrantReq) (*ListPermitGrantResp, error)
	QueryPermitGrant(c context.Context, req *QueryPermitGrantReq) (*QueryPermitGrantResp, error)
	// ListGroup
	// ListElem
	// ListGrant
}

// CreateAPIPermit CreateAPIPermit
func CreateAPIPermit(conf *config.Config) (APIPermit, error) {
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	p := &apiPermit{
		conf:      conf,
		db:        db,
		groupRepo: mysql.NewAPIPermitGroupRepo(),
		elemRepo:  mysql.NewAPIPermitElemRepo(),
		grantRepo: mysql.NewAPIPermitGrantRepo(),
		redisAPI:  redisCache,
	}

	return p, nil
}

type apiPermit struct {
	conf      *config.Config
	db        *gorm.DB
	groupRepo models.APIPermitGroupRepo
	elemRepo  models.APIPermitElemRepo
	grantRepo models.APIPermitGrantRepo
	redisAPI  models.RedisCache
}

//------------------------------------------------------------------------------

// CheckPermitReq CheckPermitReq
type CheckPermitReq struct {
	AppID  string `json:"-"`
	UserID string `json:"-"`
	KeyID  string `json:"-"`
	APIID  string `json:"-"`
	Key3ID string `json:"-"`
}

// CheckPermitResp CheckPermitResp
type CheckPermitResp struct {
}

// CheckAPIPermit CheckAPIPermit
func (p *apiPermit) CheckAPIPermit(c context.Context, req *CheckPermitReq) (*CheckPermitResp, error) {
	//TODO:
	return nil, nil
}

//------------------------------------------------------------------------------

// PermitGroupResp PermitGroupResp
type PermitGroupResp struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Access    uint   `json:"access"`
	Active    uint   `json:"active"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}

// AddPermitGroupReq AddPermitGroupReq
type AddPermitGroupReq struct {
	Owner     string `json:"-"`
	OwnerName string `json:"-"`
	Namespace string `json:"-"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
}

// AddPermitGroupResp AddPermitGroupResp
type AddPermitGroupResp = PermitGroupResp

// AddPermitGroup AddPermitGroup
func (p *apiPermit) AddPermitGroup(c context.Context, req *AddPermitGroupReq) (*AddPermitGroupResp, error) {
	pg := &models.APIPermitGroup{
		ID:        id2.WithPrefix(id2.ShortID(12), "pg_"),
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		Namespace: req.Namespace,
		Name:      req.Name,
		Title:     req.Title,
		Access:    0,
		Desc:      req.Desc,
		Active:    1,
	}

	err := p.groupRepo.Create(p.db, pg)
	if err != nil {
		return nil, err
	}

	return &AddPermitGroupResp{
		ID:        pg.ID,
		Owner:     pg.Owner,
		OwnerName: pg.OwnerName,
		Namespace: pg.Namespace,
		Name:      pg.Name,
		Title:     pg.Title,
		Desc:      pg.Desc,
		Active:    pg.Access,
		CreateAt:  pg.CreateAt,
		UpdateAt:  pg.UpdateAt,
	}, nil
}

//------------------------------------------------------------------------------

// DelPermitGroupReq DelPermitGroupReq
type DelPermitGroupReq struct {
	GroupPath string `json:"-"`
}

// DelPermitGroupResp DelPermitGroupResp
type DelPermitGroupResp struct {
	FullPath string `json:"fulPath"`
}

// DelPermitGroup DelPermitGroup
func (p *apiPermit) DelPermitGroup(c context.Context, req *DelPermitGroupReq) (*DelPermitGroupResp, error) {
	ns, name := apipath.Split(req.GroupPath)
	err := p.groupRepo.Delete(p.db, ns, name)
	return &DelPermitGroupResp{
		FullPath: req.GroupPath,
	}, err
}

//------------------------------------------------------------------------------

// UpdatePermitGroupReq updatePermitGroupReq
type UpdatePermitGroupReq struct {
	GroupPath string `json:"-"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
}

// UpdatePermitGroupResp updatePermitGroupResp
type UpdatePermitGroupResp struct {
	FullPath string `json:"fullPath"`
}

// UpdatePermitGroup UpdatePermitGroup
func (p *apiPermit) UpdatePermitGroup(c context.Context, req *UpdatePermitGroupReq) (*UpdatePermitGroupResp, error) {
	ns, name := apipath.Split(req.GroupPath)
	pg := &models.APIPermitGroup{
		Namespace: ns,
		Name:      name,
		Title:     req.Title,
		Desc:      req.Desc,
	}
	err := p.groupRepo.Update(p.db, pg)
	return &UpdatePermitGroupResp{
		FullPath: req.GroupPath,
	}, err
}

//------------------------------------------------------------------------------

// ActivePermitGroupReq ActivePermitGroupResp
type ActivePermitGroupReq struct {
	GroupPath string `json:"-"`
	Active    uint   `json:"active"`
}

// ActivePermitGroupResp ActivePermitGroupResp
type ActivePermitGroupResp struct {
	FullPath string `json:"fullPath"`
}

// ActivePermitGroup activePermitGroup
func (p *apiPermit) ActivePermitGroup(c context.Context, req *ActivePermitGroupReq) (*ActivePermitGroupResp, error) {
	ns, name := apipath.Split(req.GroupPath)
	gp := &models.APIPermitGroup{
		Namespace: ns,
		Name:      name,
		Active:    req.Active,
	}
	err := p.groupRepo.UpdateActive(p.db, gp)
	return &ActivePermitGroupResp{
		FullPath: req.GroupPath,
	}, err
}

//------------------------------------------------------------------------------

// ListPermitGroupReq ListPermitGroupReq
type ListPermitGroupReq struct {
	Namespace string `json:"-"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

// ListPermitGroupResp ListPermitGroupResp
type ListPermitGroupResp struct {
	Total int                    `json:"total"`
	Page  int                    `json:"page"`
	List  []*ListPermitGroupNode `json:"list"`
}

// ListPermitGroupNode listPermitGroupNode
type ListPermitGroupNode = PermitGroupResp

// ListPermitGroup listPermitGroup
func (p *apiPermit) ListPermitGroup(c context.Context, req *ListPermitGroupReq) (*ListPermitGroupResp, error) {
	listPG, err := p.groupRepo.List(p.db, req.Namespace, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*ListPermitGroupNode, 0, len(listPG.List))
	for _, v := range listPG.List {
		list = append(list, &ListPermitGroupNode{
			ID:        v.ID,
			Owner:     v.Owner,
			OwnerName: v.OwnerName,
			Namespace: v.Namespace,
			Name:      v.Name,
			Title:     v.Title,
			Desc:      v.Desc,
			Access:    v.Access,
			Active:    v.Active,
			CreateAt:  v.CreateAt,
			UpdateAt:  v.UpdateAt,
		})
	}

	return &ListPermitGroupResp{
		Page:  listPG.Page,
		Total: int(listPG.Total),
		List:  list,
	}, nil
}

//------------------------------------------------------------------------------

// QueryPermitGroupReq QueryPermitGroupReq
type QueryPermitGroupReq struct {
	GroupPath string `json:"-"`
}

// QueryPermitGroupResp QueryPermitGroupResp
type QueryPermitGroupResp = PermitGroupResp

// QueryPermitGroup QueryPermitGroup
func (p *apiPermit) QueryPermitGroup(c context.Context, req *QueryPermitGroupReq) (*QueryPermitGroupResp, error) {
	ns, name := apipath.Split(req.GroupPath)
	pg, err := p.groupRepo.Query(p.db, ns, name)
	if err != nil {
		return nil, err
	}
	return &QueryPermitGroupResp{
		ID:        pg.ID,
		Owner:     pg.Owner,
		OwnerName: pg.OwnerName,
		Namespace: pg.Namespace,
		Name:      name,
		Title:     pg.Title,
		Desc:      pg.Desc,
		Access:    pg.Access,
		Active:    pg.Active,
		CreateAt:  pg.CreateAt,
		UpdateAt:  pg.UpdateAt,
	}, nil
}

//------------------------------------------------------------------------------

// PermitElemResp PermitElemResp
type PermitElemResp struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	GroupPath string `json:"groupPath"`
	ElemType  string `json:"elemType"`
	ElemID    string `json:"elemID"`
	ElemPath  string `json:"elemPath"`
	Desc      string `json:"desc"`
	ElemPri   uint   `json:"elemPri"`
	Content   string `json:"content"`
	Active    uint   `json:"active"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}

// AddPermitElemReq AddPermitElemReq
type AddPermitElemReq struct {
	Owner     string `json:"-"`
	OwnerName string `json:"-"`
	GroupPath string `json:"-"`
	ElemType  string `json:"elemType"`
	ElemID    string `json:"elemID"`
	ElemPath  string `json:"elemPath"`
	Desc      string `json:"desc"`
	Content   string `json:"content"`
}

// AddPermitElemResp AddPermitElemResp
type AddPermitElemResp = PermitElemResp

// AddPermitElem AddPermitElem
func (p *apiPermit) AddPermitElem(c context.Context, req *AddPermitElemReq) (*AddPermitElemResp, error) {
	pe := &models.APIPermitElem{
		ID:        id2.WithPrefix(id2.ShortID(12), "pg_"),
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		GroupPath: req.GroupPath,
		ElemType:  req.ElemType,
		ElemID:    req.ElemID,
		ElemPath:  req.ElemPath,
		Desc:      req.Desc,
		Content:   req.Content,
		ElemPri:   0,
		Active:    1,
	}

	err := p.elemRepo.Create(p.db, pe)

	return &AddPermitElemResp{
		ID:        pe.ID,
		Owner:     pe.Owner,
		OwnerName: pe.OwnerName,
		GroupPath: pe.GroupPath,
		ElemType:  pe.ElemType,
		ElemID:    pe.ElemID,
		ElemPath:  pe.ElemPath,
		Desc:      pe.Desc,
		Content:   pe.Content,
		ElemPri:   pe.ElemPri,
		Active:    pe.Active,
		CreateAt:  pe.CreateAt,
		UpdateAt:  pe.UpdateAt,
	}, err
}

//------------------------------------------------------------------------------

// UpdatePermitElemReq UpdatePermitElemReq
type UpdatePermitElemReq struct {
	ID      string `json:"id"`
	Desc    string `json:"desc"`
	ElemPri uint   `json:"elemPri"`
}

// UpdatePermitElemResp UpdatePermitElemResp
type UpdatePermitElemResp struct {
}

// UpdatePermitElem UpdatePermitElem
func (p *apiPermit) UpdatePermitElem(c context.Context, req *UpdatePermitElemReq) (*UpdatePermitElemResp, error) {
	err := p.elemRepo.Update(p.db, &models.APIPermitElem{
		ID:      req.ID,
		Desc:    req.Desc,
		ElemPri: req.ElemPri,
	})
	return &UpdatePermitElemResp{}, err
}

//------------------------------------------------------------------------------

// DelPermitElemReq DelPermitElemReq
type DelPermitElemReq struct {
	ID string `json:"id"`
}

// DelPermitElemResp DelPermitElemResp
type DelPermitElemResp struct {
}

// DelPermitElem DelPermitElem
func (p *apiPermit) DelPermitElem(c context.Context, req *DelPermitElemReq) (*DelPermitElemResp, error) {
	err := p.elemRepo.Delete(p.db, req.ID)
	return &DelPermitElemResp{}, err
}

//------------------------------------------------------------------------------

// ActivePermitElemReq activePermitElemReq
type ActivePermitElemReq struct {
	ID     string `json:"id"`
	Active uint   `json:"active"`
}

// ActivePermitElemResp activePermitElemResp
type ActivePermitElemResp struct {
}

func (p *apiPermit) ActivePermitElem(c context.Context, req *ActivePermitElemReq) (*ActivePermitElemResp, error) {
	err := p.elemRepo.UpdateActive(p.db, &models.APIPermitElem{
		ID:     req.ID,
		Active: req.Active,
	})
	return &ActivePermitElemResp{}, err
}

//------------------------------------------------------------------------------

// ListPermitElemReq ListPermitElemReq
type ListPermitElemReq struct {
	GroupPath string `json:"-"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

// ListPermitElemResp ListPermitElemResp
type ListPermitElemResp struct {
	Total int               `json:"total"`
	Page  int               `json:"page"`
	List  []*PermitElemResp `json:"list"`
}

// ListPermitElem ListPermitElem
func (p *apiPermit) ListPermitElem(c context.Context, req *ListPermitElemReq) (*ListPermitElemResp, error) {
	listPE, err := p.elemRepo.List(p.db, req.GroupPath, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*PermitElemResp, 0, len(listPE.List))
	for _, v := range listPE.List {
		list = append(list, &PermitElemResp{
			ID:        v.ID,
			Owner:     v.Owner,
			OwnerName: v.OwnerName,
			GroupPath: v.GroupPath,
			ElemType:  v.ElemType,
			ElemID:    v.ElemID,
			ElemPath:  v.ElemPath,
			Desc:      v.Desc,
			ElemPri:   v.ElemPri,
			Content:   v.Content,
			Active:    v.Active,
			CreateAt:  v.CreateAt,
			UpdateAt:  v.UpdateAt,
		})
	}

	return &ListPermitElemResp{
		Page:  listPE.Page,
		Total: int(listPE.Total),
		List:  list,
	}, nil
}

//------------------------------------------------------------------------------

// QueryPermitElemReq QueryPermitElemReq
type QueryPermitElemReq struct {
	ID string `json:"id"`
}

// QueryPermitElemResp QueryPermitElemResp
type QueryPermitElemResp = PermitElemResp

func (p *apiPermit) QueryPermitElem(c context.Context, req *QueryPermitElemReq) (*QueryPermitElemResp, error) {
	pe, err := p.elemRepo.Query(p.db, req.ID)
	if err != nil {
		return nil, err
	}
	return &QueryPermitElemResp{
		ID:        pe.ID,
		Owner:     pe.Owner,
		OwnerName: pe.OwnerName,
		GroupPath: pe.GroupPath,
		ElemType:  pe.ElemType,
		ElemID:    pe.ElemID,
		ElemPath:  pe.ElemPath,
		Desc:      pe.Desc,
		ElemPri:   pe.ElemPri,
		Content:   pe.Content,
		Active:    pe.Active,
		CreateAt:  pe.CreateAt,
		UpdateAt:  pe.UpdateAt,
	}, nil
}

//------------------------------------------------------------------------------

// PermitGrantResp PermitGrantResp
type PermitGrantResp struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	GroupPath string `json:"groupPath"`
	GrantType string `json:"grantType"`
	GrantID   string `json:"grantID"`
	GrantName string `json:"grantName"`
	GrantPri  uint   `json:"grantPri"`
	Active    uint   `json:"active"`
	Desc      string `json:"desc"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}

// AddPermitGrantReq AddPermitGrantReq
type AddPermitGrantReq struct {
	GroupPath string `json:"-"`
	Owner     string `json:"-"`
	OwnerName string `json:"-"`
	GrantType string `json:"grantType"`
	GrantID   string `json:"grantID"`
	GrantName string `json:"grantName"`
	Desc      string `json:"desc"`
}

// AddPermitGrantResp AddPermitGrantResp
type AddPermitGrantResp = PermitGrantResp

// AddPermitGrant AddPermitGrant
func (p *apiPermit) AddPermitGrant(c context.Context, req *AddPermitGrantReq) (*AddPermitGrantResp, error) {
	pg := &models.APIPermitGrant{
		ID:        id2.WithPrefix(id2.ShortID(12), "pg_"),
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		GroupPath: req.GroupPath,
		GrantType: req.GrantType,
		GrantID:   req.GrantID,
		GrantName: req.GrantName,
		Desc:      req.Desc,
		GrantPri:  0,
		Active:    1,
	}
	err := p.grantRepo.Create(p.db, pg)
	return &AddPermitGrantResp{
		ID:        pg.ID,
		Owner:     pg.Owner,
		OwnerName: pg.OwnerName,
		GroupPath: pg.GroupPath,
		GrantType: pg.GrantType,
		GrantID:   pg.GrantID,
		GrantName: pg.GrantName,
		GrantPri:  pg.GrantPri,
		Active:    pg.Active,
		Desc:      pg.Desc,
		CreateAt:  pg.CreateAt,
		UpdateAt:  pg.UpdateAt,
	}, err
}

//------------------------------------------------------------------------------

// DelPermitGrantReq DelPermitGrantReq
type DelPermitGrantReq struct {
	ID string `json:"id"`
}

// DelPermitGrantResp DelPermitGrantResp
type DelPermitGrantResp struct {
}

// DelPermitGrant DelPermitGrant
func (p *apiPermit) DelPermitGrant(c context.Context, req *DelPermitGrantReq) (*DelPermitGrantResp, error) {
	err := p.grantRepo.Delete(p.db, req.ID)
	return &DelPermitGrantResp{}, err
}

//------------------------------------------------------------------------------

// UpdatePermitGrantReq UpdatePermitGrantReq
type UpdatePermitGrantReq struct {
	ID       string `json:"id"`
	GrantPri uint   `json:"grantPri"`
	Desc     string `json:"desc"`
}

// UpdatePermitGrantResp UpdatePermitGrantResp
type UpdatePermitGrantResp struct {
}

func (p *apiPermit) UpdatePermitGrant(c context.Context, req *UpdatePermitGrantReq) (*UpdatePermitGrantResp, error) {
	err := p.grantRepo.Update(p.db, &models.APIPermitGrant{
		ID:       req.ID,
		GrantPri: req.GrantPri,
		Desc:     req.Desc,
	})
	return &UpdatePermitGrantResp{}, err
}

//------------------------------------------------------------------------------

// ActivePermitGrantReq ActivePermitGrantReq
type ActivePermitGrantReq struct {
	ID     string `json:"id"`
	Active uint   `json:"active"`
}

// ActivePermitGrantResp ActivePermitGrantResp
type ActivePermitGrantResp struct {
}

func (p *apiPermit) ActivePermitGrant(c context.Context, req *ActivePermitGrantReq) (*ActivePermitGrantResp, error) {
	err := p.grantRepo.UpdateActive(p.db, &models.APIPermitGrant{
		ID:     req.ID,
		Active: req.Active,
	})
	return &ActivePermitGrantResp{}, err
}

//------------------------------------------------------------------------------

// ListPermitGrantReq ListPermitGrantReq
type ListPermitGrantReq struct {
	GroupPath string `json:"-"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

// ListPermitGrantResp ListPermitGrantResp
type ListPermitGrantResp struct {
	Total int                `json:"total"`
	Page  int                `json:"page"`
	List  []*PermitGrantResp `json:"list"`
}

func (p *apiPermit) ListPermitGrant(c context.Context, req *ListPermitGrantReq) (*ListPermitGrantResp, error) {
	listPG, err := p.grantRepo.List(p.db, req.GroupPath, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*PermitGrantResp, 0, len(listPG.List))
	for _, v := range listPG.List {
		list = append(list, &PermitGrantResp{
			ID:        v.ID,
			Owner:     v.Owner,
			OwnerName: v.OwnerName,
			GroupPath: v.GroupPath,
			GrantType: v.GrantType,
			GrantID:   v.GrantID,
			GrantName: v.GrantName,
			GrantPri:  v.GrantPri,
			Active:    v.Active,
			Desc:      v.Desc,
			CreateAt:  v.CreateAt,
			UpdateAt:  v.UpdateAt,
		})
	}

	return &ListPermitGrantResp{
		Page:  listPG.Page,
		Total: int(listPG.Total),
		List:  list,
	}, nil
}

//------------------------------------------------------------------------------

// QueryPermitGrantReq QueryPermitGrantReq
type QueryPermitGrantReq struct {
	ID string `json:"id"`
}

// QueryPermitGrantResp QueryPermitGrantResp
type QueryPermitGrantResp = PermitGrantResp

func (p *apiPermit) QueryPermitGrant(c context.Context, req *QueryPermitGrantReq) (*QueryPermitGrantResp, error) {
	pg, err := p.grantRepo.Query(p.db, req.ID)
	if err != nil {
		return nil, err
	}
	return &QueryPermitGrantResp{
		ID:        pg.ID,
		Owner:     pg.Owner,
		OwnerName: pg.OwnerName,
		GroupPath: pg.GroupPath,
		GrantType: pg.GrantType,
		GrantID:   pg.GrantID,
		GrantName: pg.GrantName,
		GrantPri:  pg.GrantPri,
		Active:    pg.Active,
		Desc:      pg.Desc,
		CreateAt:  pg.CreateAt,
		UpdateAt:  pg.UpdateAt,
	}, nil
}

//------------------------------------------------------------------------------
