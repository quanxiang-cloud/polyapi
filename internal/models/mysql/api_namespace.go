package mysql

import (
	"fmt"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"

	"gorm.io/gorm"
)

type apiNamespace struct {
}

// NewAPINamespaceRepo NewAPINamespaceRepo
func NewAPINamespaceRepo() models.APINamespaceRepo {
	return &apiNamespace{}
}
func (g *apiNamespace) TableName() string {
	return "api_namespace"
}

func (g *apiNamespace) updateSubCount(db *gorm.DB, fullPath string) error {
	var cnt struct {
		Parent   string
		SubCount uint
	}

	const statSubNSCountSQL = " SELECT `parent`, COUNT(1) sub_count FROM `api_namespace` WHERE `parent`=? GROUP BY `parent`; "
	const updateSubNSCountSQL = " UPDATE `api_namespace` u SET `sub_count`=? WHERE u.`parent`=? AND u.`namespace`=?; "

	if err := db.Raw(statSubNSCountSQL, fullPath).Scan(&cnt).Error; err != nil {
		return err
	}
	p, name := apipath.Split(fullPath)
	if err := db.Exec(updateSubNSCountSQL, cnt.SubCount, p, name).Error; err != nil {
		return err
	}
	return nil
}

func (g *apiNamespace) Create(db *gorm.DB, ns *models.APINamespace) error {
	now := time2.NowUnix()
	ns.CreateAt = now
	ns.UpdateAt = now
	err := db.Table(g.TableName()).Create(ns).Error
	if err == nil {
		return g.updateSubCount(db, ns.Parent)
	}
	return err
}

// inner import without updating parent count
func (g *apiNamespace) CreateInBatches(db *gorm.DB, items []*models.APINamespace) error {
	now := time2.NowUnix()
	for _, v := range items {
		v.CreateAt = now
		v.UpdateAt = now
	}
	err := db.Table(g.TableName()).CreateInBatches(items, len(items)).Error
	return err
}

func (g *apiNamespace) Delete(db *gorm.DB, path, name string) error {
	item := &models.APINamespace{
		Parent:    path,
		Namespace: name,
	}
	err := db.Table(g.TableName()).Unscoped().Where("parent=? and namespace=?", path, name).Delete(item).Error
	if err == nil {
		return g.updateSubCount(db, path)
	}
	return err
}

func (g *apiNamespace) Update(db *gorm.DB, ns *models.APINamespace) error {
	mp := map[string]interface{}{
		"title": ns.Title,
		"desc":  ns.Desc,
	}

	err := db.Table(g.TableName()).Where("parent=? and namespace=?", ns.Parent, ns.Namespace).Updates(mp).Error
	return err
}

func (g *apiNamespace) UpdateActive(db *gorm.DB, ns *models.APINamespace) error {
	mp := map[string]interface{}{
		"active": ns.Active,
	}

	err := db.Table(g.TableName()).Where("parent=? and namespace=?", ns.Parent, ns.Namespace).Updates(mp).Error
	return err
}

func (g *apiNamespace) Query(db *gorm.DB, path, name string) (*models.APINamespace, error) {
	item := &models.APINamespace{
		Parent:    path,
		Namespace: name,
	}
	err := db.Table(g.TableName()).Where("parent=? and namespace=?", path, name).Find(item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == "" {
		return nil, errNotFound
	}
	return item, nil
}

// List list namespace
// active not work if pageSize <= 0
func (g *apiNamespace) List(db *gorm.DB, path string, active, page, pageSize int) (*models.APINamespaceList, error) {
	list := &models.APINamespaceList{
		Total: -1,
	}

	listDB := db.Table(g.TableName()).Where("valid=? and parent=?", rule.Valid, path)
	if pageSize <= 0 {
		err := listDB.Find(&list.List).Error
		if err != nil {
			return nil, err
		}
		return list, nil
	}

	if active >= 0 {
		listDB.Where("active=?", active)
	}
	listDB.Count(&list.Total)
	err := listDB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list.List).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *apiNamespace) Count(db *gorm.DB, path string) (int, error) {
	var count int64
	err := db.Table(g.TableName()).Where("parent=?", path).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// parent is the path to search
// don't return parent itself
func (g *apiNamespace) Search(db *gorm.DB, parent, namespace, title string, active, page, pageSize int, withSub bool) (*models.APINamespaceList, error) {
	list := &models.APINamespaceList{}

	search := db.Table(g.TableName())
	if !withSub {
		search.Where("parent=?", parent)
	} else {
		// BUG: consider "system/foo/api" & "system/foo_001/api2" case
		search.Where("`parent`=? or `parent` like ?",
			parent, fmt.Sprintf("%s/%%", parent))
	}

	if namespace != "" {
		search.Where("namespace like ?", fmt.Sprintf("%%%s%%", namespace))
	}
	if title != "" {
		search.Where("title like ?", fmt.Sprintf("%%%s%%", title))
	}
	if active >= 0 {
		search.Where("active = ?", active)
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

func (g *apiNamespace) UpdateValidWithSub(db *gorm.DB, path string, valid uint) error {
	parent, namespace := apipath.Split(path)
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(g.TableName()).Where("`parent`=? or `parent` like ? or (`parent`=? and `namespace` = ?)",
		path, fmt.Sprintf("%s/%%", path), parent, namespace).UpdateColumn("valid", valid).Error
}

func (g *apiNamespace) DelByPrefixPath(db *gorm.DB, path string) error {
	parent, ns := apipath.Split(path)
	// BUG: consider "system/foo/api" & "system/foo_001/api2" case
	return db.Table(g.TableName()).Where("`parent`=? or `parent` like ? or (`parent`=? and `namespace`=?)",
		path, fmt.Sprintf("%s/%%", path), parent, ns).Unscoped().Delete(&models.APINamespace{}).Error
}
