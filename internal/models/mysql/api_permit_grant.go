package mysql

import (
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"

	"gorm.io/gorm"
)

type apiPermitGrant struct {
}

// NewAPIPermitGrantRepo NewAPIPermitGrantRepo
func NewAPIPermitGrantRepo() models.APIPermitGrantRepo {
	return &apiPermitGrant{}
}
func (g *apiPermitGrant) TableName() string {
	return "api_permit_grant"
}

func (g *apiPermitGrant) Create(db *gorm.DB, group *models.APIPermitGrant) error {
	now := time2.NowUnix()
	group.CreateAt = now
	group.UpdateAt = now
	return db.Table(g.TableName()).Create(group).Error
}

func (g *apiPermitGrant) Delete(db *gorm.DB, id string) error {
	item := &models.APIPermitGrant{
		ID: id,
	}
	return db.Table(g.TableName()).Unscoped().Where("id=?", id).Delete(item).Error
}

func (g *apiPermitGrant) Update(db *gorm.DB, item *models.APIPermitGrant) error {
	mp := map[string]interface{}{
		"grant_pri": item.GrantPri,
		"desc":      item.Desc,
	}

	err := db.Table(g.TableName()).Where("id=?", item.ID).Updates(mp).Error
	return err
}

func (g *apiPermitGrant) UpdateActive(db *gorm.DB, data *models.APIPermitGrant) error {
	mp := map[string]interface{}{
		"active": data.Active,
	}

	err := db.Table(g.TableName()).Where("id=?", data.ID).Updates(mp).Error
	return err
}

func (g *apiPermitGrant) Query(db *gorm.DB, id string) (*models.APIPermitGrant, error) {
	item := &models.APIPermitGrant{
		ID: id,
	}
	err := db.Table(g.TableName()).Where("id=?", id).Find(item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

func (g *apiPermitGrant) QueryGrant(db *gorm.DB, groupPath, grantType, grantID string) (*models.APIPermitGrant, error) {
	item := &models.APIPermitGrant{}
	err := db.Table(g.TableName()).Where("group_path=? and grant_type=? and grant_id=?",
		groupPath, grantType, grantID).Find(item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

func (g *apiPermitGrant) List(db *gorm.DB, groupPath string, page, pageSize int) (*models.ListPermitGrantResp, error) {
	resp := &models.ListPermitGrantResp{
		Page:  page,
		Total: -1,
	}
	listDB := db.Table(g.TableName()).Where("group_path=?", groupPath)
	if pageSize <= 0 {
		err := listDB.Find(&resp.List).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := listDB.Count(&resp.Total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&resp.List).Error
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

type grantGroups struct {
	GrantType string
	GrantID   string
	GroupPath string
}

func (g *apiPermitGrant) ListGroups(db *gorm.DB, grantType, grantID string) ([]string, error) {
	var items []*grantGroups
	err := db.Table(g.TableName()).Where("group_type=? and grant_id=?",
		grantType, grantID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errNotFound
	}
	groups := make([]string, 0, len(items))
	for _, v := range items {
		groups = append(groups, v.GroupPath)
	}
	return groups, nil
}

func (g *apiPermitGrant) Count(db *gorm.DB, groupID string) (int, error) {
	var count int64
	err := db.Table(g.TableName()).Where("group_id=?", groupID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
