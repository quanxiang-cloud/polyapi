package mysql

import (
	"fmt"
	"time"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"gorm.io/gorm"
)

type schema struct{}

// NewAPISchemaRepo NewAPISchemaRepo
func NewAPISchemaRepo() models.APISchemaRepo {
	return &schema{}
}

func (s *schema) TableName() string {
	return "api_schema"
}

func (s *schema) Create(db *gorm.DB, schema *models.APISchemaFull) error {
	if err := s.tryUpdate(db, schema); err == nil {
		return nil
	}

	t := time2.NowUnix()
	schema.CreateAt = t
	schema.UpdateAt = t
	return db.Table(s.TableName()).Create(schema).Error
}

func (s *schema) tryUpdate(db *gorm.DB, schema *models.APISchemaFull) error {
	var f struct {
		ID string
	}
	if err := db.Table(s.TableName()).Where("namespace=? and name=?", schema.Namespace, schema.Name).Find(&f).Error; err != nil {
		return err
	}
	if f.ID == "" {
		return errNotFound
	}

	err := db.Table(s.TableName()).Where("namespace=? and name=?", schema.Namespace, schema.Name).Updates(map[string]interface{}{
		"title":     schema.Title,
		"desc":      schema.Desc,
		"schema":    schema.Schema,
		"update_at": time.Now(),
	}).Error
	schema.ID = f.ID

	return err
}

func (s *schema) Delete(db *gorm.DB, namespace, name string) error {
	return db.Table(s.TableName()).Where("namespace=? and name=?", namespace, name).Delete(&models.APISchemaFull{}).Error
}

func (s *schema) Query(db *gorm.DB, namespace, name string) (*models.APISchemaFull, error) {
	ret := new(models.APISchemaFull)
	err := db.Table(s.TableName()).Where("namespace=? and name=?", namespace, name).Find(ret).Error
	return ret, err
}

func (s *schema) List(db *gorm.DB, namespace string, withSub bool, page, pageSize int) (*models.APISchemaList, error) {
	ret := &models.APISchemaList{
		Total: -1,
	}

	sql := db.Table(s.TableName())
	if withSub {
		sql = sql.Where("namespace like ?", fmt.Sprintf("%s%%", namespace))
	} else {
		sql = sql.Where("namespace=?", fmt.Sprintf("%s%%", namespace))
	}

	if pageSize <= 0 {
		if err := sql.Find(&ret.List).Error; err != nil {
			return nil, err
		}
		return ret, nil
	}

	if err := sql.Count(&ret.Total).Error; err != nil {
		return nil, err
	}

	if err := sql.Offset((page - 1) * pageSize).Limit(pageSize).Find(&ret.List).Error; err != nil {
		return nil, err
	}
	return ret, nil
}
