package mysql

import (
	"fmt"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"

	"gorm.io/gorm"
)

type raw struct {
}

// TableName TableName
func (r *raw) TableName() string {
	return "api_raw"
}

// NewRawAPIRepo NewRawAPIRepo
func NewRawAPIRepo() models.RawAPIRepo {
	return &raw{}

}

func (r *raw) Create(db *gorm.DB, raw *models.RawAPIFull) error {
	// try update first
	if err := r.tryUpdate(db, raw); err == nil {
		return nil
	}

	now := time2.NowUnix()
	raw.CreateAt = now
	raw.UpdateAt = now
	return db.Table(r.TableName()).Create(raw).Error
}

func (r *raw) CreateInBatches(db *gorm.DB, items []*models.RawAPIFull) error {
	now := time2.NowUnix()
	for _, item := range items {
		item.CreateAt = now
		item.UpdateAt = now
	}
	return db.Table(r.TableName()).CreateInBatches(items, len(items)).Error
}

func (r *raw) tryUpdate(db *gorm.DB, raw *models.RawAPIFull) error {
	var f struct {
		ID string
	}
	if err := db.Table(r.TableName()).Where("namespace=? and name=?",
		raw.Namespace, raw.Name).Find(&f).Error; err != nil {
		return err
	}
	if f.ID == "" {
		return errNotFound
	}

	raw.UpdateAt = time2.NowUnix()
	err := db.Table(r.TableName()).Where("ID=?", f.ID).Updates(map[string]interface{}{
		"title":     raw.Title,
		"path":      raw.Path,
		"url":       raw.URL,
		"action":    raw.Action,
		"method":    raw.Method,
		"version":   raw.Version,
		"desc":      raw.Desc,
		"content":   raw.Content,
		"doc":       raw.Doc,
		"update_at": raw.UpdateAt,
	}).Error
	raw.ID = f.ID
	return err
}

func (r *raw) Del(db *gorm.DB, namespace string, names []string) error {
	return db.Table(r.TableName()).Where("namespace=? and name in (?)", namespace, names).Unscoped().Delete(&models.RawAPIFull{}).Error
}

func (r *raw) Get(db *gorm.DB, path, name string) (*models.RawAPICore, error) {
	raw := new(models.RawAPICore)
	err := db.Table(r.TableName()).Where("namespace=? and name=?", path, name).Find(raw).Error
	if err != nil {
		return nil, err
	}
	if raw.ID == "" {
		return nil, errNotFound
	}
	return raw, nil
}

func (r *raw) GetInBatches(db *gorm.DB, path [][2]string) (*models.RawAPIList, error) {
	list := make([]*models.RawAPICore, 0)
	err := r.getInBatches(db, path, &list)
	return &models.RawAPIList{
		List: list,
	}, err
}

// func (r *raw) GetFullInBatches(db *gorm.DB, path [][2]string) ([]*models.RawAPIFull, error) {
// 	ret := make([]*models.RawAPIFull, 0)
// 	err := r.getInBatches(db, path, &ret)
// 	return ret, err
// }

func (r *raw) GetDocInBatches(db *gorm.DB, path [][2]string) ([]*models.RawAPIDoc, error) {
	ret := make([]*models.RawAPIDoc, 0)
	err := r.getInBatches(db, path, &ret)
	return ret, err
}

func (r *raw) getInBatches(db *gorm.DB, path [][2]string, ret interface{}) error {
	return db.Table(r.TableName()).Where("(namespace, name) in ?", path).Find(ret).Error
}

func (r *raw) GetDoc(db *gorm.DB, path, name string) (*models.RawAPIDoc, error) {
	raw := new(models.RawAPIDoc)
	err := db.Table(r.TableName()).Where("namespace=? and name=?", path, name).Find(raw).Error
	if err != nil {
		return nil, err
	}
	if raw.ID == "" {
		return nil, errNotFound
	}
	return raw, nil
}

// GetByID get the raw api
func (r *raw) GetByID(db *gorm.DB, id string) (*models.RawAPICore, error) {
	raw := new(models.RawAPICore)
	err := db.Table(r.TableName()).Where("id=?", id).Find(raw).Error
	if err != nil {
		return nil, err
	}
	if raw.ID == "" {
		return nil, errNotFound
	}
	return raw, nil
}

// List list raw by page, if pageSize <= 0, list all.
// active not work when pageSize <= 0
func (r *raw) List(db *gorm.DB, namespace, service string, active, page, pageSize int) (*models.RawAPIList, error) {
	var list = &models.RawAPIList{
		Total: -1,
	}

	listDB := db.Table(r.TableName()).Where("delete_at is null and valid=?", rule.Valid)
	if service != "" {
		listDB = listDB.Where("service=?", service)
	} else {
		listDB = listDB.Where("namespace=?", namespace)
	}

	if pageSize <= 0 {
		if err := listDB.Find(&list.List).Error; err != nil {
			return nil, err
		}
		return list, nil
	}

	if active >= 0 {
		listDB.Where("active = ?", active)
	}
	listDB.Count(&list.Total)
	if err := listDB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list.List).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *raw) ListByPrefixPath(db *gorm.DB, path string, active, page, pageSize int) ([]*models.RawAPIFull, int64, error) {
	list := make([]*models.RawAPIFull, 0)
	var total int64 = -1
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	listDB := db.Table(r.TableName()).Where("(`namespace`=? or `namespace` like ?) and delete_at is null",
		path, fmt.Sprintf("%s/%%", path))
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

func (r *raw) UpdateActive(db *gorm.DB, namespace, name string, active uint) error {
	return db.Table(r.TableName()).Where("namespace=? and name=?", namespace, name).UpdateColumn("active", active).Error
}

func (r *raw) UpdateValid(db *gorm.DB, path [][2]string, valid uint) error {
	return db.Table(r.TableName()).Where("(namespace, name) in ?", path).UpdateColumn("valid", valid).Error
}

func (r *raw) UpdateInBatch(db *gorm.DB, namespace, service, host, schema, authType string) error {
	mp := map[string]interface{}{
		"host":      host,
		"schema":    schema,
		"auth_type": authType,
		"url":       gorm.Expr("concat(?, path)", fmt.Sprintf("%s://%s", schema, host)),
	}
	return db.Table(r.TableName()).Where("service=?", service).Updates(mp).Error
}

func (r *raw) Search(db *gorm.DB, namespace, name, title string, active, page, pageSize int, withSub bool) (*models.RawAPIList, error) {
	list := &models.RawAPIList{}

	search := db.Table(r.TableName())
	if !withSub {
		search.Where("namespace=?", namespace)
	} else {
		// BUG: consider "system/foo/api" & "system/foo_001/api2" case
		//      search "system/foo/api" don't list out "system/foo_001/xxx"
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
		search.Where("active=?", active)
	}

	search.Where("valid=?", rule.Valid)

	if pageSize <= 0 {
		err := search.Find(&list.List).Error
		if err != nil {
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

func (r *raw) UpdateValidByPrefixPath(db *gorm.DB, path string, valid uint) error {
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(r.TableName()).Where("`namespace`=? or `namespace` like ?",
		valid, fmt.Sprintf("%s/%%", path)).UpdateColumn("valid", valid).Error
}

func (r *raw) DelByPrefixPath(db *gorm.DB, path string) error {
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(r.TableName()).Where("`namespace`=? or `namespace` like ?",
		path, fmt.Sprintf("%s/%%", path)).Unscoped().Delete(&models.RawAPICore{}).Error
}
