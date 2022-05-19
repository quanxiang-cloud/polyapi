package expr

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
)

// flexFactory regist some flex json objects
var flexFactory = factory.NewFlexObjFactory("flexObjCreator")

// GetFactory return the factory object
func GetFactory() *factory.FlexObjFactory { return flexFactory }

func init() {
	flexFactory.MustReg(ValNumber("123"))
	flexFactory.MustReg(ValString("foo"))
	flexFactory.MustReg(ValBoolean("true"))
	flexFactory.MustReg(ValObject{
		Value{
			Name: "a",
			Type: "string",
			Data: FlexJSONObject{
				D: NewStringer("foo", ValTypeString),
			},
		},
		Value{
			Name: "b",
			Type: "number",
			Data: FlexJSONObject{
				D: NewStringer("123", ValTypeNumber),
			},
		},
		Value{
			Name: "c",
			Type: "boolean",
			Data: FlexJSONObject{
				D: NewStringer("true", ValTypeBoolean),
			},
		},
	})
	flexFactory.MustReg(ValArray{
		Value{
			Type: "string",
			Data: FlexJSONObject{
				D: ValString("foo"),
			},
		},
		Value{
			Type: "string",
			Data: FlexJSONObject{
				D: ValString("bar"),
			},
		},
	})
	flexFactory.MustReg(ValArrayString{"foo", "bar"})

	flexFactory.MustReg(ValTimestamp("2020-12-31T12:34:56Z"))
	flexFactory.MustReg(ValAction("createUser"))
}
