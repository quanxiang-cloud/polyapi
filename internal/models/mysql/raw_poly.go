package mysql

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/internal/models"

	"gorm.io/gorm"
)

// NewRawPolyRepo NewRawPolyRepo
func NewRawPolyRepo() models.RawPolyRepo {
	return &rawPoly{}
}

type rawPoly struct {
}

func (rp *rawPoly) TableName() string {
	return "api_raw_poly"
}

func (rp *rawPoly) Create(db *gorm.DB, item *models.RawPoly) error {
	return db.Table(rp.TableName()).Create(item).Error
}

func (rp *rawPoly) CreateInBatches(db *gorm.DB, items []*models.RawPoly) error {
	return db.Table(rp.TableName()).CreateInBatches(items, len(items)).Error
}

func (rp *rawPoly) DeleteByRawAPI(db *gorm.DB, rawPath string) error {
	return db.Table(rp.TableName()).Where("raw_api=?", rawPath).Delete(&models.RawPoly{}).Error
}

func (rp *rawPoly) DeleteByPolyAPI(db *gorm.DB, polyPath string) error {
	return db.Table(rp.TableName()).Where("poly_api=?", polyPath).Delete(&models.RawPoly{}).Error
}

func (rp *rawPoly) DeleteByPolyAPIInBatches(db *gorm.DB, polyPath []string) error {
	return db.Table(rp.TableName()).Where("poly_api in (?)", polyPath).Delete(&models.RawPoly{}).Error
}

func (rp *rawPoly) QueryByRawAPI(db *gorm.DB, rawPath []string) (*models.RawPolyList, error) {
	var total int64
	var list []*models.RawPoly
	query := db.Table(rp.TableName()).Where("raw_api in (?)", rawPath)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := query.Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return &models.RawPolyList{
		List:  list,
		Total: total,
	}, nil
}

func (rp *rawPoly) QueryByPolyAPI(db *gorm.DB, polyPath string) (*models.RawPolyList, error) {
	var total int64
	var list []*models.RawPoly
	query := db.Table(rp.TableName()).Where("poly_api=?", polyPath)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		if err := query.Find(&list).Error; err != nil {
			return nil, err
		}
	}
	return &models.RawPolyList{
		List:  list,
		Total: total,
	}, nil
}

func (rp *rawPoly) Update(db *gorm.DB, items []*models.RawPoly) error {
	// BUG: crash when list is empty
	if len(items) == 0 {
		return nil
	}
	// TODO: to remove not in items and add not exists
	tx := db.Begin()
	if err := rp.DeleteByPolyAPI(tx, items[0].PolyAPI); err != nil {
		tx.Rollback()
		return err
	}
	if err := rp.CreateInBatches(tx, items); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (rp *rawPoly) ListByPrefixPath(db *gorm.DB, path string) (*models.RawPolyList, error) {
	ret := &models.RawPolyList{
		Total: -1,
	}
	err := db.Table(rp.TableName()).Where("raw_api like ?", fmt.Sprintf("%s%%", path)).Find(&ret.List).Error
	return ret, err
}

func (rp *rawPoly) DelByPrefixPath(db *gorm.DB, path string) error {
	condition := fmt.Sprintf("%s%%", path)
	err := db.Table(rp.TableName()).Where("raw_api like ? or poly_api like ?", condition, condition).Unscoped().Delete(&models.RawPoly{}).Error
	return err
}
