package rawpoly

import (
	"fmt"
	"strings"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobcenter"
	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"

	"gorm.io/gorm"
)

func init() {
	inst.reg("apiPath", "update api path if has change", updateAPIPath)

	jobcenter.RegRunner(300, jobName, fmt.Sprintf("*Update table %s\n%s", tb, inst.show()), &inst)
}

//------------------------------------------------------------------------------

const (
	jobName   = "rawpoly"
	tb        = "api_raw_poly"
	batchSize = 200
)

type dbRecord = adaptor.RawPoly

type updateFunc func(data *dbRecord) (bool, error)
type handler struct {
	name string
	desc string
	fn   updateFunc
}
type job struct {
	handlers []*handler
}

var inst job

func (j *job) reg(name string, desc string, fn updateFunc) {
	j.handlers = append(j.handlers, &handler{name, desc, fn})
}

func (j *job) show() string {
	buf := strings.Builder{}
	for _, v := range j.handlers {
		buf.WriteString(fmt.Sprintf("  -%-10s %s\n", v.name, v.desc))
	}
	return buf.String()
}

func (j *job) update(data *dbRecord) (updated bool, err error) {
	for _, v := range j.handlers {
		u, err := v.fn(data)
		if err != nil {
			err = fmt.Errorf("%s:%s", v.name, err.Error())
		}
		updated = updated || u
	}
	return
}

// update api_raw.doc from swagger
func (j job) Run(db *gorm.DB) (status string, err error) {
	defer func() {
		d := db.Exec("DELETE FROM `api_raw_poly` WHERE `poly_api` NOT IN (SELECT CONCAT(`namespace`,'/',`name`)FROM `api_poly`) OR `raw_api` NOT IN (SELECT CONCAT(`namespace`,'/',`name`)FROM `api_raw`) ").RowsAffected
		if d > 0 {
			status += fmt.Sprintf(" d=%d", d)
		}
	}()

	if renamed.Raw.Empty() && renamed.Poly.Empty() {
		return "no need to update", nil
	}

	t := time.Now()
	last := t

	var result []*dbRecord
	total := 0
	errCnt := 0
	updateCnt := 0
	batchProcess := func(tx *gorm.DB, batch int) error {
		var update = make([]*dbRecord, 0, len(result))
		for _, v := range result {
			switch updated, err := j.update(v); {
			case err != nil:
				errCnt++
				fmt.Printf("**%sErr cnt=%d poly=%q and raw=%q err=%v\n", jobName, errCnt, v.PolyAPI, v.RawAPI, err)

				fallthrough
			default:
				if updated {
					updateCnt++
					update = append(update, v)
				}
			}
		}

		if len(update) > 0 {
			if err := db.Table(tb).Save(update).Error; err != nil {
				return err
			}
		}

		total += len(result)
		fmt.Printf("==%s batch=%d total=%d\n", jobcenter.Now(t, &last), batch, total)
		return nil
	}

	err = db.Table(tb).Where("1=?", 1).FindInBatches(&result, batchSize, batchProcess).Error
	return fmt.Sprintf("e=%d,u=%d,t=%d", errCnt, updateCnt, total), err
}
