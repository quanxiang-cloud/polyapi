package jsvm

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/protocol"

	"github.com/quanxiang-cloud/cabin/logger"
)

const (
	debug = consts.Debug // debug mode

	// MaxVMNum max vm count
	MaxVMNum         = 32
	waitTooManyTimes = 3                // weakup but get vm failed times limit
	vmHeaderVarName  = "__pdHttpHeader" // http header variable in vm
)

// regexTrimWhitespace remove whitespace from JSON
var regexTrimWhitespace = regexp.MustCompile(`(?s-m:\s+)`)

// js vms manager
var mgr = newManager()

// NewEvalerCreator make a evalerCreator
func NewEvalerCreator() *EvalerCreator {
	return &EvalerCreator{}
}

// EvalerCreator object
type EvalerCreator struct {
}

// CreateEvaler create a protocol.Evaler object
func (c *EvalerCreator) CreateEvaler() protocol.Evaler {
	return CreateVM()
}

// CreateVM create a vm object.
// NOTE: don't forget call vm.Free() after use to release the vm.
//       Or vm will leak and block other vm request.
func CreateVM() *JsVM {
	return mgr.GetVM()
}

// EvalScript check if js script is error
func EvalScript(script string) error {
	vm := mgr.GetVM()
	defer func() {
		vm.Free()
	}()

	return vm.Eval(script)
}

// RunJsString run a js script in a free vm
// It will wait a free vm if vms are all busy
func RunJsString(script string, inputJSONBody []byte, header http.Header) (string, error) {
	start := time.Now()
	vm := mgr.GetVM()
	if debug {
		fmt.Println("getVM", vm.id)
	}
	vm.SetHTTPHeader(header)
	defer func() {
		vm.Free()
		if debug {
			dur := time.Now().Sub(start)
			e := len(script)
			if e > 10 {
				e = 10
			}
			fmt.Printf("****run script cost %s in vm %d [%#v]\n", dur, vm.id, script[:e]+"...")
		}
	}()
	if len(inputJSONBody) == 0 {
		inputJSONBody = []byte("{}")
	}

	/*
		var __input = {
			header: __pdHttpHeader,
			body: pdFromJson(inputJSONBody),
		 x : {...},
		}
	*/

	// remove whitespace from JSON
	singleLinebody := unsafeByteString(regexTrimWhitespace.ReplaceAll(inputJSONBody, []byte("")))
	const varName = consts.PolyAPIInputVarName

	polyScript := fmt.Sprintf(`
// prepare input:
var %s = {
	header: %s,
	body: pdFromJson('%s'),
	x: {
	},
}
// function call:
%s
`, varName, vmHeaderVarName, singleLinebody, script)

	//println("run script:", polyScript) // TODO: remove

	v, err := vm.RunString(polyScript)
	if err != nil {
		logger.Logger.Errorf("[jsvm] execute fail: %s \nscript='%s'", err.Error(), script)
		return "", errcode.ErrVMExecuteFail.NewError()
	}

	out, ok := exportValue(v).(string)
	if !ok {
		return "", errors.New("script doesnt return string")
	}
	return out, nil
}

// JsVM hold a js vm object
type JsVM struct {
	*vmRuntime
	id     int64
	mgr    *VmsManager
	header http.Header
}

// Free makes the vm free for other goroutine
func (vm *JsVM) Free() {
	vm.mgr.Free(vm)
}

// SetHTTPHeader set httpHeader variable in vm
func (vm *JsVM) SetHTTPHeader(header http.Header) {
	vm.header = header
	if vm.header == nil {
		vm.header = make(http.Header)
	}
	vm.setVMHTTPHeader()
}

// AddHTTPHeader add httpHeader variable in vm (export to vm)
func (vm *JsVM) AddHTTPHeader(key, value string) {
	if vm.header == nil {
		vm.header = make(http.Header)
	}
	vm.header.Add(key, value)
	vm.setVMHTTPHeader()
}

// NewHTTPHeader create a http header
func (vm *JsVM) NewHTTPHeader() http.Header {
	h := http.Header{}
	if refer := vm.header.Get(consts.HeaderRefer); refer != "" {
		h.Set(consts.HeaderRefer, refer) //Copy header Refer
	}
	return h
}

// QueryUser return user-id in header
func (vm *JsVM) QueryUser(realDo bool) string {
	if realDo {
		return auth.GetAuthOwner(vm.header)
	}
	return ""
}

// QueryUserName return user-name in header
func (vm *JsVM) QueryUserName() string {
	return auth.GetAuthOwnerName(vm.header)
}

func (vm *JsVM) setVMHTTPHeader() {
	vm.Set(vmHeaderVarName, vm.header)
}

func newManager() *VmsManager {
	p := &VmsManager{
		free: make([]*JsVM, 0, MaxVMNum),
	}
	p.cond = sync.NewCond(&p.lock)
	return p
}

// VmsManager manage a multi set of js vm objects
// multi vms request by multi goroutines
// each one to serve one goroutine
type VmsManager struct {
	cond  *sync.Cond // wait list of vm request
	lock  sync.Mutex // lock of free list
	free  []*JsVM    // list of free vms
	idGen int64      // atomic access
}

func (m *VmsManager) newJsVM(id int64) *JsVM {
	vm := newSysVM()
	p := &JsVM{
		vmRuntime: vm,
		id:        id,
		mgr:       m,
	}
	p.vmSetup()
	return p
}

// GetVM get a free JsVm for one goroutine
func (m *VmsManager) GetVM() *JsVM {
	m.lock.Lock()
	defer m.lock.Unlock()
	for i := 0; ; i++ {
		if n := len(m.free); n > 0 {
			vm := m.free[n-1]
			m.free = m.free[:n-1]
			return vm
		}
		if atomic.LoadInt64(&m.idGen) <= MaxVMNum {
			return m.newJsVM(atomic.AddInt64(&m.idGen, 1))
		}
		m.cond.Wait() // wait for a free action to weakup me
		if debug {
			fmt.Println("weakup for GetVM", i)
		}

		// it is expected to get vm success if weakup by free,
		// but it's probably to hijack by new httpHandler
		// hijack waitTooManyTimes(3) times means system too busy
		if i > waitTooManyTimes {
			// TODO: no panic
			panic(fmt.Errorf("GetVM wait too many times %d", i))
		}
	}

	return nil
}

// Free put the vm into free list
func (m *VmsManager) Free(vm *JsVM) {
	if vm.mgr != m {
		panic("Free on unexpected vm")
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.free = append(m.free, vm)

	if debug {
		fmt.Println("freeVM", vm.id)
	}

	m.cond.Signal() // notify wait list
}
