package jsvm

import (
	"fmt"
	"testing"

	pkgtimestamp "github.com/quanxiang-cloud/polyapi/pkg/basic/timestamp"
)

func TestPredefined(t *testing.T) {
	type testCase struct {
		code, expect string
	}
	cases := []*testCase{
		&testCase{
			`"123"`,
			`123`,
		},
		&testCase{
			`timestamp("foo")`,
			pkgtimestamp.Timestamp(""),
		},
		&testCase{
			`pdSelect(2<3, "sel-true", "sel-false")`,
			`sel-true`,
		},
		&testCase{
			`pdSelect(2>3, "sel-true", "sel-false")`,
			`sel-false`,
		},
		&testCase{
			`pdToJsonP(pdMergeObjs(
				{a:1, b:2, },
				{c:3, d:4, }
			))`,
			`{
  "a": 1,
  "b": 2,
  "c": 3,
  "d": 4
}`,
		},
		&testCase{
			`pdToJsonP(pdFiltObject(
				{a:1, b:2, c:3},
				{black:{b:"",}, white:{b:"bb",}},
				function(i, d){ return undefined }
			))`,
			`{
  "a": 1,
  "bb": 2,
  "c": 3
}`,
		},
		&testCase{
			`pdToJsonP(pdFiltObject(
				[
				 	{a:1, b:2, c:3},
					{a:2, b:3, c:4},
					{a:3, b:4, c:5},
					{a:4, b:5, c:6},
				],
				{black:{b:"",}, white:{c:"cc",}},
				function(i, d){ return d.b>=3&&d.b<=4 }
			))`,
			`[
  {
    "a": 2,
    "cc": 4
  },
  {
    "a": 3,
    "cc": 5
  }
]`,
		},
		&testCase{
			`pdToJsonP(pdFiltObject(
				[
				 	{a:1, b:2, c:3},
					{a:2, b:3, c:4},
					{a:3, b:4, c:5},
					{a:4, b:5, c:6},
				],
				{black:{}, white:{c:"cc",}},
				function(i, d){ return d.b==3||d.b==5 }
			))`,
			`[
  {
    "cc": 4
  },
  {
    "cc": 6
  }
]`,
		},
		// &testCase{
		// 	`timestamp()`,
		// 	``,
		// },
		&testCase{
			`
	var x = "123"
	"pre${x}pro"
			`,
			`pre${x}pro`,
		},
		&testCase{
			`""`,
			``,
		},
		&testCase{
			`
	var x = "123"
	format("|%s|%s", x, "456")
			`,
			`|123|456`,
		},
		// &testCase{
		// 	`
		// var d = {
		// a: 123,
		// b: "xyz",
		// c: [1, "x", true],
		// d: {
		// x: "xx",
		// y: 345,
		// z: true
		// },
		// }
		// pdToJsonP(pdToJsobj("xml", pdToXml(d)))
		// 	`,
		// 	`{
		// "d": {
		//   "x": "xx",
		//   "y": 345,
		//   "z": true
		// },
		// "a": 123,
		// "b": "xyz",
		// "c": [
		//   1,
		//   "x",
		//   true
		// ]
		// }`,
		// },
		// &testCase{
		// 	`
		// var d = {
		// a: 123,
		// b: "xyz",
		// c: [1, "x", true],
		// d: {
		// x: "xx",
		// y: 345,
		// z: true
		// },
		// }
		// //pdToJsonP(pdToJsobj("yaml", pdToYaml(d)))
		// pdToXml(pdToJsobj("yaml", pdToYaml(d)))
		// 	`,
		// 	``,
		// },
		&testCase{
			`
var x = '{"code":90014020001,"msg":"permission group name is exist","data":null}'
pdToJsonP(pdToJsobj("json", x))
			`,
			`{
  "code": 90014020001,
  "msg": "permission group name is exist",
  "data": null
}`,
		},
		&testCase{
			`
var _tmp = function(){
  var d = { "__input": __input, } // qyAllLocalData

  d.start = __input.body

  d.start.header = d.start.header || {}

  if (true) { // req1, create
    var _apiPath = "http://fanyi.youdao.com/translate"
    var _t = {
      "doctype": "json",
      "type": "AUTO",
      "i": "计算",
    }
    var _th = pdNewHttpHeader()
    pdAddHttpHeader(_th, "Content-Type", "application/json")

    var _tk = '';
    var _tb = pdAppendAuth(_tk, 'none', _th, pdToJson(_t))
    d.req1 = pdToJsobj("json", pdHttpRequest(_apiPath, "GET", _tb, _th, pdQueryUser(true)))
  }

  d.end = {
    "req1": d.req1.translateResult,
  }
  return pdToJsonP(d.end)
}; _tmp();
			`,
			`{
  "req1": [
    [
      {
        "src": "计算",
        "tgt": "To calculate"
      }
    ]
  ]
}`,
		},
		&testCase{
			`eval('"a"+"b"')`,
			`ab`,
		},
	}
	for i, v := range cases {
		got, err := RunJsString(v.code, []byte(`{"a":1, "b":2}`), nil)
		if got != v.expect || err != nil {
			t.Errorf("case %d, expect [%q] got [%q], error: %v", i, v.expect, got, err)
			fmt.Println(got)
		}
	}
}

func TestEval(t *testing.T) {
	type testCase struct {
		code        string
		expectError bool
	}
	cases := []*testCase{
		&testCase{`a+b`, false},
		&testCase{`a.x+b.y`, false},
		&testCase{`887`, false},
		&testCase{`"abc"`, false},
		&testCase{`abc`, false},
		&testCase{`8f56eda`, true},
		&testCase{`f8f56eda`, false},
		&testCase{`a.x@b.y`, true},
	}
	vm := CreateVM()
	defer vm.Free()
	for i, v := range cases {
		err := vm.Eval(v.code)
		if got := err != nil; got != v.expectError {
			t.Errorf("case %d, expect [%v] got [%v], error: %v", i, v.expectError, got, err)
			fmt.Println(i, err)
		}
	}
}

func TestJson(t *testing.T) {
	fmt.Printf(`{"$":%q}`, "abc\n")
}
