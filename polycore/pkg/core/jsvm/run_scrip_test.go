package jsvm_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/stat"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/jsvm"
)

const (
	testDebug = false
)

var (
	beginNum int32 = 0
	endNum   int32 = 0
	wg       sync.WaitGroup
	st       = stat.NewTimeStat("RunJsScript")
)

func _TestRunJsScript(t *testing.T) {
	code := `
	function __qyTmpFn(){
		//pdHttpAddHeader("Content-Type", "application/json")
		//pdHttpAddHeader("Access-Token", "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5")
		var _th = pdNewHttpHeader()
		var x=qyHttpPost("http://keeper.test/api/v1/org/getUserTemplate", JSON.stringify({}), _th)
		var y=qyHttpPost("http://keeper.test/api/v1/structor/group/list", JSON.stringify({appID:"appIDbalabala"}), _th)
		var xx=JSON.parse(x)
		var yy=JSON.parse(y)
		var t={code:xx.code,data:{count:yy.data.count+100,fileURL:xx.data.fileURL+"_addSuffix"}}
		return JSON.stringify(t, undefined, 2)
	};__qyTmpFn();
`
	code2 := `
	(function (){
		//pdHttpAddHeader("Content-Type", "application/json")
		//pdHttpAddHeader("Access-Token", "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5")
		var _th = pdNewHttpHeader()
		var x=qyHttpPost("http://keeper.test/api/v1/org/getUserTemplate", qyToJson({}), _th)
		var y=qyHttpPost("http://keeper.test/api/v1/structor/group/list", qyToJson({appID:"appIDbalabala"}), _th)
		var xx=qyFromJson(x)
		var yy=qyFromJson(y)
		var t={code:xx.code,data:{count:yy.data.count+111,fileURL:xx.data.fileURL+"_addSuffix2"}}
		t["xx-yy"]="xxx-yyy"
		return qyToJsonP(t)
	}());
`
	code3 := `
	var _tmp = function(){

var d = { "__input": __input, } // qyAllLocalData

d.start = __input.body

d.start.header = d.start.header || {}

if (true) { // DYqIdUzn,

var _apiPath = format("https://home.yunify.com:443/distributor.action" )

var _t = {

"serviceName": "clogin",

"userName": d.start.userName ,

"password": d.start.password ,

}

var _th = pdNewHttpHeader()

//pdAddHttpHeader(_th, "Content-Type", "application/json")

var _tk = '';

var _tb = pdAppendAuth(_tk, 'none', _th, pdToJson(_t))

d.DYqIdUzn = pdToJsobj("json", pdHttpRequest(_apiPath, "GET", _tb, _th, pdQueryUser(true)))

}

if (true) { // pryLZHnL,

var _apiPath = format("https://home.yunify.com:443/distributor.action" )

var _t = {

"serviceName": d.start.serviceName ,

"objectApiName": d.start.objectApiName ,

//"expressions": "ownerid='"+d.DYqIdUzn.userInfo.userId+"'",

"binding": d.DYqIdUzn.binding ,

}

var _th = pdNewHttpHeader()

//pdAddHttpHeader(_th, "Content-Type", "application/json")

var _tk = '';

var _tb = pdAppendAuth(_tk, 'none', _th, pdToJson(_t))

d.pryLZHnL = pdToJsobj("json", pdHttpRequest(_apiPath, "GET", _tb, _th, pdQueryUser(true)))

}

d.end = {

"result": d.pryLZHnL.result,

"data": d.pryLZHnL.data,

}

return pdToJsonP(d.end)

}; _tmp();
`
	batch := jsvm.MaxVMNum*0 + 1
	for i := 0; i < batch; i++ {
		wg.Add(1)

		if groupWait := time.Millisecond * 500; groupWait > 0 {
			if group := jsvm.MaxVMNum / 2; (i % group) == (group - 1) {
				time.Sleep(groupWait)
			}
		}
		if false {
			go testRunJsScript(i+1, code, false)
		}

		if false {
			go testRunJsScript(i+1, code2, true)
		}

		if true {
			testRunJsScript(i+1, code3, true)
		}
	}
	wg.Wait()
	if testDebug {
		fmt.Println("time stat:", st.Report())
	}
}

func testRunJsScript(id int, code string, gotErr bool) {
	s := time.Now()
	if testDebug {
		atomic.AddInt32(&beginNum, 1)
		fmt.Println(id, atomic.LoadInt32(&beginNum), "begin")
	}

	v, err := jsvm.RunJsString(code, nil, nil)
	if got := err != nil; got != gotErr {
		fmt.Println(v, err)
	}

	if testDebug {
		atomic.AddInt32(&endNum, 1)
		fmt.Println(id, atomic.LoadInt32(&endNum), "finish")
	}
	st.Add(time.Now().Sub(s))
	wg.Done()
}
