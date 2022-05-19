package rawapi

import (
	"fmt"
	"strings"
	"time"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobcenter"

	"gorm.io/gorm"
)

func init() {
	inst.reg("name", "uuid api name 'raw_xxx' => short name", updateName)
	inst.reg("suffix", "update name from 'foo' => 'foo.r'", updateSuffix)
	inst.reg("doc", "update 'doc' from 'doc.swagger'", updateDoc)

	jobcenter.RegRunner(100, jobName, fmt.Sprintf("*Update table api_raw\n%s", inst.show()), &inst)
}

//------------------------------------------------------------------------------

const (
	jobName   = "rawapi"
	tb        = "api_raw"
	batchSize = 50
)

type dbRecord = models.RawAPIFull

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
func (j job) Run(db *gorm.DB) (string, error) {
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
				fmt.Printf("**%sErr cnt=%d namespace=%q and name=%q err=%v\n", jobName, errCnt, v.Namespace, v.Name, err)

				fallthrough
			default:
				if updated {
					updateCnt++
					v.UpdateAt = time2.NowUnix()
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

	err := db.Table(tb).Where("1=?", 1).FindInBatches(&result, batchSize, batchProcess).Error
	return fmt.Sprintf("e=%d,u=%d,t=%d", errCnt, updateCnt, total), err
}
