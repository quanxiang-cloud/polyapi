package mysql

import (
	"fmt"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"

	"gorm.io/gorm"
)

type apiService struct {
}

// NewAPIServiceRepo NewAPIServiceRepo
func NewAPIServiceRepo() models.APIServiceRepo {
	return &apiService{}
}

func (g *apiService) TableName() string {
	return "api_service"
}

func (g *apiService) Create(db *gorm.DB, group *models.APIService) error {
	now := time2.NowUnix()
	group.CreateAt = now
	group.UpdateAt = now
	return db.Table(g.TableName()).Create(group).Error
}

func (g *apiService) CreateInBatches(db *gorm.DB, groups []*models.APIService) error {
	now := time2.NowUnix()
	for _, v := range groups {
		v.CreateAt = now
		v.UpdateAt = now
	}
	return db.Table(g.TableName()).CreateInBatches(groups, len(groups)).Error
}

func (g *apiService) Delete(db *gorm.DB, namespace, name string) error {
	item := &models.APIService{}
	return db.Table(g.TableName()).Unscoped().
		Where("namespace=? and name=?", namespace, name).Delete(item).Error
}

func (g *apiService) DeleteInBatch(db *gorm.DB, namespace string, names []string) error {
	return db.Table(g.TableName()).Where("namespace=? and name in (?)", namespace, names).Delete(&models.APIService{}).Error
}

func (g *apiService) Update(db *gorm.DB, group *models.APIService) error {
	mp := map[string]interface{}{
		"title": group.Title,
		"desc":  group.Desc,
	}

	err := db.Table(g.TableName()).Where("namespace=? and name=?", group.Namespace, group.Name).Updates(mp).Error
	return err
}

func (g *apiService) UpdateProperty(db *gorm.DB, group *models.APIService) error {
	mp := map[string]interface{}{
		"schema":    group.Schema,
		"host":      group.Host,
		"auth_type": group.AuthType,
		"authorize": group.Authorize,
	}
	return db.Table(g.TableName()).Where("namespace=? and name=?", group.Namespace, group.Name).Updates(mp).Error
}

func (g *apiService) UpdateActive(db *gorm.DB, group *models.APIService) error {
	mp := map[string]interface{}{
		"active": group.Active,
	}

	err := db.Table(g.TableName()).Where("namespace=? and name=?", group.Namespace, group.Name).Updates(mp).Error
	return err
}

func (g *apiService) Query(db *gorm.DB, namespace, name string) (*models.APIService, error) {
	item := &models.APIService{}
	err := db.Table(g.TableName()).
		Where("namespace=? and name=?", namespace, name).Find(item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

func (g *apiService) List(db *gorm.DB, path string, page, pageSize int, withSub bool) (*models.APIServiceList, error) {
	var list = &models.APIServiceList{
		Total: -1,
	}

	listDB := db.Table(g.TableName())
	if !withSub {
		listDB = listDB.Where("namespace=?", path)
	} else {
		listDB = listDB.Where("namespace like ?", fmt.Sprintf("%s%%", path))
	}

	if pageSize <= 0 {
		err := listDB.Find(&list.List).Error
		if err != nil {
			return nil, err
		}
		return list, nil
	}

	listDB.Count(&list.Total)
	err := listDB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list.List).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *apiService) Count(db *gorm.DB, path string) (int, error) {
	var count int64
	err := db.Table(g.TableName()).Where("namespace=?", path).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (g *apiService) DelByPrefixPath(db *gorm.DB, path string) error {
	return db.Table(g.TableName()).Where("namespace like ?", fmt.Sprintf("%s%%", path)).Unscoped().Delete(&models.APIService{}).Error
}
