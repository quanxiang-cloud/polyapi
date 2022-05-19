package mysql

import (
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"

	"gorm.io/gorm"
)

type apiPermitElem struct {
}

// NewAPIPermitElemRepo NewAPIPermitElemRepo
func NewAPIPermitElemRepo() models.APIPermitElemRepo {
	return &apiPermitElem{}
}
func (g *apiPermitElem) TableName() string {
	return "api_permit_elem"
}

func (g *apiPermitElem) Create(db *gorm.DB, group *models.APIPermitElem) error {
	now := time2.NowUnix()
	group.CreateAt = now
	group.UpdateAt = now
	return db.Table(g.TableName()).Create(group).Error
}

func (g *apiPermitElem) Delete(db *gorm.DB, id string) error {
	item := &models.APIPermitElem{
		ID: id,
	}
	return db.Table(g.TableName()).Unscoped().Where("id=?", id).Delete(item).Error
}

func (g *apiPermitElem) Update(db *gorm.DB, item *models.APIPermitElem) error {
	mp := map[string]interface{}{
		"elem_pri": item.ElemPri,
		"desc":     item.Desc,
	}

	err := db.Table(g.TableName()).Where("id=?", item.ID).Updates(mp).Error
	return err
}

func (g *apiPermitElem) UpdateActive(db *gorm.DB, data *models.APIPermitElem) error {
	mp := map[string]interface{}{
		"active": data.Active,
	}

	err := db.Table(g.TableName()).Where("id=?", data.ID).Updates(mp).Error
	return err
}

func (g *apiPermitElem) Query(db *gorm.DB, id string) (*models.APIPermitElem, error) {
	item := &models.APIPermitElem{
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

func (g *apiPermitElem) QueryElem(db *gorm.DB, groups []string, elemType, elemID string) (*models.APIPermitElem, error) {
	item := &models.APIPermitElem{}
	err := db.Table(g.TableName()).Where("groupPath in ? and elem_type=? and elem_id=?",
		groups, elemType, elemID).Find(item).Limit(1).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

func (g *apiPermitElem) List(db *gorm.DB, groupPath string, page, pageSize int) (*models.ListPermitElemResp, error) {
	resp := &models.ListPermitElemResp{
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

func (g *apiPermitElem) Count(db *gorm.DB, groupID string) (int, error) {
	var count int64
	err := db.Table(g.TableName()).Where("group_id=?", groupID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
