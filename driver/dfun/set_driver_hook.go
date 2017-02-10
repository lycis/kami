package dfun

import (
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// set_driver_hook is used to set specific driver callback functions.
//
// This functions requires the executing script to have ROOT privileges.
//
// Hooks:
//   H_HB_ERROR_FUNC - function is called when an error occurs in the
//                     heartbeat call

const (
	H_HB_ERROR_FUNC = iota
)

type dfun_set_driver_hook struct {
	script *script.ScriptContext
}

func create_dfun_set_driver_hook(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_set_driver_hook{
		script: i,
	}
}

func (e dfun_set_driver_hook) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

func (df dfun_set_driver_hook) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		return otto.UndefinedValue()
	}
}
