package mysql

import (
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"

	"gorm.io/gorm"
)

type apiPermitGroup struct {
}

// NewAPIPermitGroupRepo NewAPIPermitGroupRepo
func NewAPIPermitGroupRepo() models.APIPermitGroupRepo {
	return &apiPermitGroup{}
}
func (g *apiPermitGroup) TableName() string {
	return "api_permit_group"
}

func (g *apiPermitGroup) Create(db *gorm.DB, group *models.APIPermitGroup) error {
	now := time2.NowUnix()
	group.CreateAt = now
	group.UpdateAt = now
	return db.Table(g.TableName()).Create(group).Error
}

func (g *apiPermitGroup) Delete(db *gorm.DB, path, name string) error {
	item := &models.APIPermitGroup{
		Namespace: path,
		Name:      name,
	}
	return db.Table(g.TableName()).Unscoped().Where("namespace=? and name=?", path, name).Delete(item).Error
}

func (g *apiPermitGroup) Update(db *gorm.DB, item *models.APIPermitGroup) error {
	mp := map[string]interface{}{
		"title": item.Title,
		"desc":  item.Desc,
	}

	err := db.Table(g.TableName()).Where("namespace=? and name=?", item.Namespace, item.Name).Updates(mp).Error
	return err
}

func (g *apiPermitGroup) UpdateActive(db *gorm.DB, item *models.APIPermitGroup) error {
	mp := map[string]interface{}{
		"active": item.Active,
	}

	err := db.Table(g.TableName()).Where("namespace=? and name=?",
		item.Namespace, item.Name).Updates(mp).Error
	return err
}

func (g *apiPermitGroup) Query(db *gorm.DB, namespacePath, name string) (*models.APIPermitGroup, error) {
	item := &models.APIPermitGroup{
		Namespace: namespacePath,
		Name:      name,
	}
	err := db.Table(g.TableName()).Where("namespace=? and name=?",
		namespacePath, name).Find(item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

func (g *apiPermitGroup) List(db *gorm.DB, path string, page, pageSize int) (*models.ListPermitGroupResp, error) {
	resp := &models.ListPermitGroupResp{
		Page:  page,
		Total: -1,
	}
	listDB := db.Table(g.TableName()).Where("namespace=?", path)
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

func (g *apiPermitGroup) Count(db *gorm.DB, path string) (int, error) {
	var count int64
	err := db.Table(g.TableName()).Where("namespace=?", path).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
