package mysql

import (
	"fmt"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"

	"gorm.io/gorm"
)

// NewPolyDbAPI NewPolyDbAPI
func NewPolyDbAPI() models.PolyAPIRepo {
	return &polyDb{}
}

type polyDb struct {
}

// TableName of poly api db
func (r *polyDb) TableName() string {
	return models.PolyTableName
}

// Create create a new poly api
func (r *polyDb) Create(db *gorm.DB, info *models.PolyAPIArrange) error {
	now := time2.NowUnix()
	info.CreateAt = now
	info.UpdateAt = now
	info.BuildAt = now
	//d.DeleteAt = now
	return db.Model(info).Create(info).Error
}

func (r *polyDb) CreateInBatches(db *gorm.DB, items []*models.PolyAPIFull) error {
	now := time2.NowUnix()
	for _, item := range items {
		item.CreateAt = now
		item.UpdateAt = now
		item.BuildAt = now
	}
	return db.Table(r.TableName()).CreateInBatches(items, len(items)).Error
}

// Delete delete a poly api
func (r *polyDb) Delete(db *gorm.DB, path string, name []string) error {
	return db.Table(r.TableName()).Where("namespace=? and name in (?)", path, name).Unscoped().Delete(&models.PolyAPIArrange{}).Error
}

// UpdateArrange update a poly api arrange json
func (r *polyDb) UpdateArrange(db *gorm.DB, info *models.PolyAPIArrange) error {
	now := time2.NowUnix()
	info.UpdateAt = now
	err := db.Table(r.TableName()).
		Where("namespace=? and name=?", info.Namespace, info.Name).Updates(info).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateScript save the build result
func (r *polyDb) UpdateScript(db *gorm.DB, info *models.PolyBuildResult) error {
	now := time2.NowUnix()
	info.BuildAt = now
	err := db.Table(r.TableName()).
		Where("namespace=? and name=?", info.Namespace, info.Name).Updates(info).Error
	if err != nil {
		return err
	}
	return nil
}

// GetArrange return a poly api arrange json
func (r *polyDb) GetArrange(db *gorm.DB, path, name string) (*models.PolyAPIArrange, error) {
	d := &models.PolyAPIArrange{}
	err := db.Table(r.TableName()).Where("namespace=? and name=?", path, name).Find(d).Error
	if err != nil {
		return nil, err
	}
	if d.ID == "" {
		return nil, errNotFound
	}
	return d, nil
}

// GetScript get the build result
func (r *polyDb) GetScript(db *gorm.DB, path, name string) (*models.PolyAPIScript, error) {
	d := &models.PolyAPIScript{}
	err := db.Table(r.TableName()).Where("namespace=? and name=?", path, name).Find(d).Error
	if err != nil {
		return nil, err
	}
	if d.ID == "" {
		return nil, errNotFound
	}
	return d, nil
}

// GetDoc get the doc info
func (r *polyDb) GetDoc(db *gorm.DB, path, name string) (*models.PolyAPIDoc, error) {
	d := &models.PolyAPIDoc{}
	err := db.Table(r.TableName()).Where("namespace=? and name=?", path, name).Find(d).Error
	if err != nil {
		return nil, err
	}
	if d.ID == "" {
		return nil, errNotFound
	}
	return d, nil
}

func (r *polyDb) GetDocInBatches(db *gorm.DB, path [][2]string) ([]*models.PolyAPIDoc, error) {
	ret := make([]*models.PolyAPIDoc, 0)
	err := db.Table(r.TableName()).Where("(namespace, name) in (?)", path).Find(&ret).Error
	return ret, err
}

// List list poly api
// active not work when pageSize <= 0
func (r *polyDb) List(db *gorm.DB, namespace string, active, page, pageSize int) (*models.PolyAPIList, error) {
	var list = &models.PolyAPIList{
		Total: -1,
	}

	listDB := db.Table(r.TableName()).Where("namespace=? and delete_at is null and valid=?", namespace, rule.Valid)

	if pageSize <= 0 {
		if err := listDB.Find(&list.List).Error; err != nil {
			return nil, err
		}
		return list, nil
	}

	if active >= 0 {
		listDB.Where("active=?", active)
	}
	listDB.Count(&list.Total)
	if err := listDB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list.List).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *polyDb) UpdateActive(db *gorm.DB, namespace, name string, active uint) error {
	return db.Table(r.TableName()).Where("namespace=? and name=?", namespace, name).UpdateColumn("active", active).Error
}

func (r *polyDb) Search(db *gorm.DB, namespace, name, title string, active, page, pageSize int, withSub bool) (*models.PolyAPIList, error) {
	list := &models.PolyAPIList{}

	search := db.Table(r.TableName())
	if !withSub {
		search.Where("namespace=?", namespace)
	} else {
		// BUG: consider "system/foo/api" & "system/foo_001/api2" case
		search.Where("`namespace`=? or `namespace` like ?",
			namespace, fmt.Sprintf("%s/%%", namespace))
	}

	if name != "" {
		search.Where("name like ?", fmt.Sprintf("%%%s%%", name))
	}
	if title != "" {
		search.Where("title like ?", fmt.Sprintf("%%%s%%", title))
	}
	if active >= 0 {
		search.Where("active = ?", active)
	}

	search.Where("valid=?", rule.Valid)

	if pageSize <= 0 {
		if err := search.Find(&list.List).Error; err != nil {
			return nil, err
		}
		list.Total = int64(len(list.List))
		return list, nil
	}

	err := search.Count(&list.Total).Error
	if err != nil {
		return nil, err
	}
	err = search.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list.List).Error
	return list, err
}

func (r *polyDb) UpdateValid(db *gorm.DB, path [][2]string, valid uint) error {
	return db.Table(r.TableName()).Where("(namespace, name) in ?", path).UpdateColumn("valid", valid).Error
}

func (r *polyDb) UpdateValidByPrefixPath(db *gorm.DB, namespace string, valid uint) error {
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(r.TableName()).Where("`namespace`=? or `namespace` like ?",
		namespace, fmt.Sprintf("%s/%%", namespace)).UpdateColumn("valid", valid).Error
}

func (r *polyDb) ListByPrefixPath(db *gorm.DB, namespace string, active, page, pageSize int) ([]*models.PolyAPIFull, int64, error) {
	list := make([]*models.PolyAPIFull, 0)
	var total int64 = -1
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	listDB := db.Table(r.TableName()).Where("(`namespace`=? or `namespace` like ?) and delete_at is null",
		namespace, fmt.Sprintf("%s/%%", namespace))
	if active >= 0 {
		listDB.Where("active=?", active)
	}

	if pageSize <= 0 {
		if err := listDB.Find(&list).Error; err != nil {
			return nil, total, err
		}
		return list, total, nil
	}

	listDB.Count(&total)
	if err := listDB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}

func (r *polyDb) DelByPrefixPath(db *gorm.DB, path string) error {
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(r.TableName()).Where("`namespace`=? or `namespace` like ?",
		path, fmt.Sprintf("%s/%%", path)).Unscoped().Delete(&models.PolyAPIArrange{}).Error
}
