package permission

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

// permit define
var (
	// PermitEnum represents permit objects
	PermitEnum    = enumset.New(nil)
	PermitRead    = PermitEnum.MustReg("read")
	PermitExecute = PermitEnum.MustReg("execute")
	PermitCreate  = PermitEnum.MustReg("create")
	PermitUpdate  = PermitEnum.MustReg("update")
	PermitDelete  = PermitEnum.MustReg("delete")
	PermitGrant   = PermitEnum.MustReg("grant")
)

var permitName = map[string]PermitBit{}

var permitBit = map[PermitBit]string{
	PermitBitRead:    PermitRead.String(),
	PermitBitExecute: PermitExecute.String(),
	PermitBitCreate:  PermitCreate.String(),
	PermitBitUpdate:  PermitUpdate.String(),
	PermitBitDelete:  PermitDelete.String(),
	PermitBitGrant:   PermitGrant.String(),
}

func init() {
	for k, v := range permitBit {
		permitName[v] = k
	}
}

// ParsePermits parse permit bit list to bits
func ParsePermits(permits []string) (PermitBit, error) {
	var r PermitBit
	var err error
	for _, v := range permits {
		if b, ok := permitName[v]; ok {
			r |= b
		} else {
			if err == nil {
				err = fmt.Errorf(`invalid permit "%s", valid list: %v`, v, PermitEnum.GetAll())
			}
		}
	}
	return r, err
}
