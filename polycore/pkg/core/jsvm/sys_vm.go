package jsvm

import (
	"fmt"
	"strings"

	jsvm "github.com/dop251/goja"
	//jsvm "github.com/robertkrimen/otto"
)

type sysVM = jsvm.Runtime

type sysVMValue = jsvm.Value

var emptySysVMValue sysVMValue

type vmRuntime struct {
	*sysVM
}

// RunString runs string script
func (rt *JsVM) RunString(script string) (sysVMValue, error) {
	return rt.sysVM.RunString(script)
}

// Eval execute script by eval
func (rt *JsVM) Eval(script string) error {
	scriptRun := fmt.Sprintf("eval('%s')", script)
	_, err := rt.sysVM.RunString(scriptRun)
	if err != nil {
		// ignore ReferenceError:
		if strings.HasPrefix(err.Error(), "ReferenceError:") {
			return nil
		}
	}
	return err
}

// Set set Go symbol to vm
func (rt *JsVM) Set(name string, value interface{}) error {
	return rt.sysVM.Set(name, value)
}

func exportValue(v sysVMValue) interface{} {
	return v.Export()
	// d, err := v.Export()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return d
}

func newSysVM() *vmRuntime {
	return &vmRuntime{jsvm.New()}
}
